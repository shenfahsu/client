// @flow
import HiddenString from '../util/hidden-string'
import {call, put, select, fork} from 'redux-saga/effects'
import {deviceDeviceHistoryListRpcPromise, loginDeprovisionRpcPromise, loginPaperKeyRpcChannelMap, revokeRevokeDeviceRpcPromise, rekeyGetRevokeWarningRpcPromise} from '../constants/types/flow-types'
import {devicesTab, loginTab} from '../constants/tabs'
import {isMobile} from '../constants/platform'
import {is} from 'immutable'
import {navigateTo} from './route-tree'
import {safeTakeEvery, safeTakeLatest, singleFixedChannelConfig, closeChannelMap, takeFromChannelMap, effectOnChannelMap} from '../util/saga'
import {setRevokedSelf} from './login'

import type {DeviceRemoved, GeneratePaperKey, IncomingDisplayPaperKeyPhrase, LoadDevices, LoadingDevices, PaperKeyLoaded, PaperKeyLoading, RemoveDevice, ShowDevices, ShowRemovePage} from '../constants/devices'
import type {Device} from '../constants/types/more'
import type {SagaGenerator} from '../constants/types/saga'

isMobile && module.hot && module.hot.accept(() => {
  console.log('accepted update in actions/devices')
})

export function loadDevices (): LoadDevices {
  return {payload: undefined, type: 'devices:loadDevices'}
}

export function loadingDevices (): LoadingDevices {
  return {payload: undefined, type: 'devices:loadingDevices'}
}

export function removeDevice (deviceID: string, name: string, currentDevice: boolean): RemoveDevice {
  return {payload: {currentDevice, deviceID, name}, type: 'devices:removeDevice'}
}

export function showRemovePage (device: Device): ShowRemovePage {
  return {payload: {device}, type: 'devices:showRemovePage'}
}

export function generatePaperKey (): GeneratePaperKey {
  return {payload: undefined, type: 'devices:generatePaperKey'}
}

function * _deviceShowRemovePageSaga (showRemovePageAction: ShowRemovePage): SagaGenerator<any, any> {
  const device = showRemovePageAction.payload.device
  let endangeredTLFs = {endangeredTLFs: []}
  try {
    endangeredTLFs = yield call(rekeyGetRevokeWarningRpcPromise, {param: {targetDevice: device.deviceID}})
  } catch (e) {
    console.warn('Error getting endangered TLFs:', e)
  }
  yield put(navigateTo([devicesTab,
    {props: {device, endangeredTLFs}, selected: 'devicePage'},
    {props: {device, endangeredTLFs}, selected: 'removeDevice'},
  ]))
}

function * _deviceListSaga (): SagaGenerator<any, any> {
  yield put(loadingDevices())
  try {
    const devices = yield call(deviceDeviceHistoryListRpcPromise)
    yield put(({
      payload: devices,
      type: 'devices:showDevices',
    }: ShowDevices))
  } catch (e) {
    yield put(({
      error: true,
      payload: {errorObj: e, errorText: e.desc + e.name},
      type: 'devices:showDevices',
    }: ShowDevices))
  }
}

function * _deviceRemoveSaga (removeAction: RemoveDevice): SagaGenerator<any, any> {
  // Record our current route, only navigate away later if it's unchanged.
  const beforeRouteState = yield select(state => state.routeTree.routeState)

  // Revoking the current device uses the "deprovision" RPC instead.
  const {currentDevice, name, deviceID} = removeAction.payload
  if (currentDevice) {
    try {
      const username = yield select(state => state.config && state.config.username)
      if (!username) {
        const error = {errorText: 'No username in removeDevice'}
        console.warn(error)
        yield put(({
          error: true,
          payload: error,
          type: 'devices:deviceRemoved',
        }: DeviceRemoved))
      }
      yield call(loginDeprovisionRpcPromise, {param: {doRevoke: true, username}})
      yield put(navigateTo([loginTab]))
      yield put(setRevokedSelf(name))
      yield put(({
        payload: undefined,
        type: 'devices:deviceRemoved',
      }: DeviceRemoved))
    } catch (e) {
      console.warn('Error removing the current device:', e)
      yield put(({
        error: true,
        payload: {errorObj: e, errorText: e.desc + e.name},
        type: 'devices:deviceRemoved',
      }: DeviceRemoved))
    }
  } else {
    // Not the current device.
    try {
      yield call(revokeRevokeDeviceRpcPromise, {
        param: {deviceID, force: false},
      })
      yield put(({
        payload: undefined,
        type: 'devices:deviceRemoved',
      }: DeviceRemoved))
    } catch (e) {
      console.warn('Error removing a device:', e)
      yield put(({
        error: true,
        payload: {errorObj: e, errorText: e.desc + e.name},
        type: 'devices:deviceRemoved',
      }: DeviceRemoved))
    }
  }
  yield put(loadDevices())

  const afterRouteState = yield select(state => state.routeTree.routeState)
  if (is(beforeRouteState, afterRouteState)) {
    yield put(navigateTo([devicesTab]))
  }
}

function _generatePaperKey (channelConfig) {
  return loginPaperKeyRpcChannelMap(channelConfig, {})
}

function * _handlePromptRevokePaperKeys (chanMap): SagaGenerator<any, any> {
  yield effectOnChannelMap(c => safeTakeEvery(c, ({response}) => response.result(false)), chanMap, 'keybase.1.loginUi.promptRevokePaperKeys')
}

function * _devicePaperKeySaga (): SagaGenerator<any, any> {
  yield put(({
    payload: undefined,
    type: 'devices:paperKeyLoading',
  }: PaperKeyLoading))

  const channelConfig = singleFixedChannelConfig(['keybase.1.loginUi.promptRevokePaperKeys', 'keybase.1.loginUi.displayPaperKeyPhrase'])

  const generatePaperKeyChanMap = ((yield call(_generatePaperKey, channelConfig)): any)
  try {
    yield fork(_handlePromptRevokePaperKeys, generatePaperKeyChanMap)
    const displayPaperKeyPhrase:
      ?IncomingDisplayPaperKeyPhrase = ((yield takeFromChannelMap(generatePaperKeyChanMap, 'keybase.1.loginUi.displayPaperKeyPhrase')): any)
    if (!displayPaperKeyPhrase) {
      const error = {errorText: 'no displayPaperKeyPhrase response'}
      console.warn(error.errorText)
      yield put(({
        error: true,
        payload: error,
        type: 'devices:paperKeyLoaded',
      }: PaperKeyLoaded))
      return
    }
    yield put(({
      payload: new HiddenString(displayPaperKeyPhrase.params.phrase),
      type: 'devices:paperKeyLoaded',
    }: PaperKeyLoaded))
    displayPaperKeyPhrase.response.result()
  } catch (e) {
    closeChannelMap(generatePaperKeyChanMap)
    console.warn('error in generating paper key', e)
    yield put(({
      error: true,
      payload: {errorObj: e, errorText: e.desc + e.name},
      type: 'devices:paperKeyLoaded',
    }: PaperKeyLoaded))
  }
}

function * deviceSaga (): SagaGenerator<any, any> {
  yield [
    safeTakeLatest('devices:loadDevices', _deviceListSaga),
    safeTakeEvery('devices:removeDevice', _deviceRemoveSaga),
    safeTakeEvery('devices:generatePaperKey', _devicePaperKeySaga),
    safeTakeEvery('devices:showRemovePage', _deviceShowRemovePageSaga),
  ]
}

export default deviceSaga
