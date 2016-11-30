def getCommit(path) {
    dir(path) {
        return bat(returnStdout: true, script: '@echo off && git rev-parse HEAD').trim()
    }
}

def checkout_keybase(repo, revision) {
    dir("src/github.com/keybase/${repo}") {
        retry(3) {
            checkout([
                scm: [
                    $class: 'GitSCM',
                    branches: [[name: revision]],
                    userRemoteConfigs: [[url: "https://github.com/keybase/${repo}.git"]],
                ]
            ])
        }
    }
}

def publish(bucket, excluded, source) {
    step([
        $class: 'S3BucketPublisher',
        dontWaitForConcurrentBuildCompletion: false,
        entries: [[
            bucket: bucket,
            excludedFile: excluded,
            flatten: true,
            noUploadOnFailure: true,
            selectedRegion: 'us-east-1',
            sourceFile: source,
            storageClass: 'STANDARD',
            uploadFromSlave: true,
        ]],
        profileName: 'keybase',
        userMetadata: []
    ])

}

def doBuild() {
    stage('Checkout Client') {
        // Reset symlink due to node/git/windows problems
        bat 'if EXIST src\\github.com\\keybase\\client\\shared cd src\\github.com\\keybase\\client && git checkout shared'
        bat 'if EXIST src\\github.com\\keybase\\client\\desktop\\shared cd src\\github.com\\keybase\\client && git checkout desktop/shared'
        bat 'if EXIST src\\github.com\\keybase\\client\\desktop\\renderer\\fonts cd src\\github.com\\keybase\\client && git checkout desktop/renderer/fonts'
        parallel(
            checkout_client: { checkout_keybase('client', ClientRevision) },
            checkout_kbfs: { checkout_keybase('kbfs', KBFSRevision) },
            checkout_updater: { checkout_keybase('go-updater', UpdaterRevision) },
            checkout_release: { checkout_keybase('release', ReleaseRevision) },
        )
    }
    // Make sure any previous desktop build is deleted
    bat '''
        if EXIST src\\github.com\\keybase\\client\\desktop\\release rmdir /q /s src\\github.com\\keybase\\client\\desktop\\release
        path
    '''                
    stage('Build Client') {
        bat '"%ProgramFiles(x86)%\\Microsoft Visual Studio 14.0\\vc\\bin\\vcvars32.bat" && src\\github.com\\keybase\\client\\packaging\\windows\\build_prerelease.cmd'
    } 
    stage('Build UI') {
        bat 'src\\github.com\\keybase\\client\\packaging\\windows\\buildui.bat'
    }
    stage('Build Installer') {
        bat 'call "%ProgramFiles(x86)%\\Microsoft Visual Studio 14.0\\vc\\bin\\vcvars32.bat" && src\\github.com\\keybase\\client\\packaging\\windows\\doinstaller_wix.cmd'
        archiveArtifacts "src\\github.com\\keybase\\client\\packaging\\windows\\${BUILD_TAG}\\*.*"
    }

    stage('Publish to S3') {
        if (UpdateChannel != "None"){    
            publish("prerelease.keybase.io/windows", 
                    "", 
                    "src\\github.com\\keybase\\client\\packaging\\windows\\${BUILD_TAG}\\*.exe") 
            publish("prerelease.keybase.io", 
                    "", 
                    "src\\github.com\\keybase\\client\\packaging\\windows\\${BUILD_TAG}\\update-windows-prod-test-v2.json") 
        } else {
            echo "No update channel"
        }
    }

    
    stage('Invoke SmokeB build') {
        if (UpdateChannel == "Smoke"){
            def clientCommit = getCommit('src\\github.com\\keybase\\client')
            def kbfsCommit =  getCommit('src\\github.com\\keybase\\kbfs')
            def updaterCommit =  getCommit('src\\github.com\\keybase\\go-updater')
            def releaseCommit =  getCommit('src\\github.com\\keybase\\release')
            def smokeASemVer = ''
            dir('src\\github.com\\keybase\\client\\go\\keybase') {
                smokeASemVer = bat(returnStdout: true, script: '@echo off && winresource.exe -cv').trim()
            }
            build([
                job: "${env.JOB_NAME}",
                parameters: [
                    string(name: 'ClientRevision', value: "${clientCommit}"),
                    string(name: 'KBFSRevision', value: "${kbfsCommit}"),
                    string(name: 'UpdaterRevision', value: "${updaterCommit}"),
                    string(name: 'ReleaseRevision', value: "${releaseCommit}"),
                    string(name: 'DOKAN_PATH', value: "${DOKAN_PATH}"),
                    string(name: 'UpdateChannel', value: 'Smoke2'),
                    string(name: 'SmokeASemVer', value: "${smokeASemVer}"
                )],
                wait: false
            ])
        } else {
            echo "Not a smoke build"
        }
    }
    stage('Publish smoke updater jsons to S3') {
        if (UpdateChannel == "Smoke2") {
            publish("prerelease.keybase.io/windows-support", 
                "src\\github.com\\keybase\\client\\packaging\\windows\\${BUILD_TAG}\\update-windows-prod-test-v2.json",
                "src\\github.com\\keybase\\client\\packaging\\windows\\${BUILD_TAG}\\*.json") 
            def smokeBSemVer = ''
            dir('src\\github.com\\keybase\\client\\go\\keybase') {
                smokeBSemVer = bat(returnStdout: true, script: '@echo off && winresource.exe -cv').trim()
            }
            withCredentials([[
                $class: 'StringBinding',
                credentialsId: 'KEYBASE_TOKEN',
                variable: 'KEYBASE_TOKEN'
                ]]) {
                dir('src\\github.com\\keybase\\release') {
                    bat "release announce-build --build-a=\"${params.SmokeASemVer}\" --build-b=\"${smokeBSemVer}\" --platform=\"windows\""
                }
            }
        } else {
            echo "Non smoke build"
        }
    }
}

// Invoke the build with a separate workspace for each executor,
// and with GOPATH set to that workspace
node ('windows-release') {
    ws("${WORKSPACE}_${EXECUTOR_NUMBER}") {
        withEnv(["GOPATH=${pwd()}"]) {
            doBuild()
        }
    }
}