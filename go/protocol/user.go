// Auto-generated by avdl-compiler v1.3.1 (https://github.com/keybase/node-avdl-compiler)
//   Input file: avdl/user.avdl

package keybase1

import (
	rpc "github.com/keybase/go-framed-msgpack-rpc"
	context "golang.org/x/net/context"
)

type Tracker struct {
	Tracker UID  `codec:"tracker" json:"tracker"`
	Status  int  `codec:"status" json:"status"`
	MTime   Time `codec:"mTime" json:"mTime"`
}

type TrackProof struct {
	ProofType string `codec:"proofType" json:"proofType"`
	ProofName string `codec:"proofName" json:"proofName"`
	IdString  string `codec:"idString" json:"idString"`
}

type WebProof struct {
	Hostname  string   `codec:"hostname" json:"hostname"`
	Protocols []string `codec:"protocols" json:"protocols"`
}

type Proofs struct {
	Social     []TrackProof `codec:"social" json:"social"`
	Web        []WebProof   `codec:"web" json:"web"`
	PublicKeys []PublicKey  `codec:"publicKeys" json:"publicKeys"`
}

type UserSummary struct {
	Uid          UID    `codec:"uid" json:"uid"`
	Username     string `codec:"username" json:"username"`
	Thumbnail    string `codec:"thumbnail" json:"thumbnail"`
	IdVersion    int    `codec:"idVersion" json:"idVersion"`
	FullName     string `codec:"fullName" json:"fullName"`
	Bio          string `codec:"bio" json:"bio"`
	Proofs       Proofs `codec:"proofs" json:"proofs"`
	SigIDDisplay string `codec:"sigIDDisplay" json:"sigIDDisplay"`
	TrackTime    Time   `codec:"trackTime" json:"trackTime"`
}

type SearchComponent struct {
	Key   string  `codec:"key" json:"key"`
	Value string  `codec:"value" json:"value"`
	Score float64 `codec:"score" json:"score"`
}

type SearchResult struct {
	Uid        UID               `codec:"uid" json:"uid"`
	Username   string            `codec:"username" json:"username"`
	Components []SearchComponent `codec:"components" json:"components"`
	Score      float64           `codec:"score" json:"score"`
}

type UserSummary2 struct {
	Uid        UID    `codec:"uid" json:"uid"`
	Username   string `codec:"username" json:"username"`
	Thumbnail  string `codec:"thumbnail" json:"thumbnail"`
	FullName   string `codec:"fullName" json:"fullName"`
	IsFollower bool   `codec:"isFollower" json:"isFollower"`
	IsFollowee bool   `codec:"isFollowee" json:"isFollowee"`
}

type UserSummary2Set struct {
	Users   []UserSummary2 `codec:"users" json:"users"`
	Time    Time           `codec:"time" json:"time"`
	Version int            `codec:"version" json:"version"`
}

type ListTrackersArg struct {
	SessionID int `codec:"sessionID" json:"sessionID"`
	Uid       UID `codec:"uid" json:"uid"`
}

type ListTrackersByNameArg struct {
	SessionID int    `codec:"sessionID" json:"sessionID"`
	Username  string `codec:"username" json:"username"`
}

type ListTrackersSelfArg struct {
	SessionID int `codec:"sessionID" json:"sessionID"`
}

type LoadUncheckedUserSummariesArg struct {
	SessionID int   `codec:"sessionID" json:"sessionID"`
	Uids      []UID `codec:"uids" json:"uids"`
}

type LoadUserArg struct {
	SessionID int `codec:"sessionID" json:"sessionID"`
	Uid       UID `codec:"uid" json:"uid"`
}

type LoadUserByNameArg struct {
	SessionID int    `codec:"sessionID" json:"sessionID"`
	Username  string `codec:"username" json:"username"`
}

type LoadUserPlusKeysArg struct {
	SessionID int `codec:"sessionID" json:"sessionID"`
	Uid       UID `codec:"uid" json:"uid"`
}

type LoadPublicKeysArg struct {
	SessionID int `codec:"sessionID" json:"sessionID"`
	Uid       UID `codec:"uid" json:"uid"`
}

type ListTrackingArg struct {
	SessionID int    `codec:"sessionID" json:"sessionID"`
	Filter    string `codec:"filter" json:"filter"`
	Assertion string `codec:"assertion" json:"assertion"`
}

type ListTrackingJSONArg struct {
	SessionID int    `codec:"sessionID" json:"sessionID"`
	Filter    string `codec:"filter" json:"filter"`
	Verbose   bool   `codec:"verbose" json:"verbose"`
	Assertion string `codec:"assertion" json:"assertion"`
}

type SearchArg struct {
	SessionID int    `codec:"sessionID" json:"sessionID"`
	Query     string `codec:"query" json:"query"`
}

type LoadAllPublicKeysUnverifiedArg struct {
	SessionID int `codec:"sessionID" json:"sessionID"`
	Uid       UID `codec:"uid" json:"uid"`
}

type ListTrackers2Arg struct {
	SessionID int    `codec:"sessionID" json:"sessionID"`
	Assertion string `codec:"assertion" json:"assertion"`
	Reverse   bool   `codec:"reverse" json:"reverse"`
}

type UserInterface interface {
	ListTrackers(context.Context, ListTrackersArg) ([]Tracker, error)
	ListTrackersByName(context.Context, ListTrackersByNameArg) ([]Tracker, error)
	ListTrackersSelf(context.Context, int) ([]Tracker, error)
	// Load user summaries for the supplied uids.
	// They are "unchecked" in that the client is not verifying the info from the server.
	// If len(uids) > 500, the first 500 will be returned.
	LoadUncheckedUserSummaries(context.Context, LoadUncheckedUserSummariesArg) ([]UserSummary, error)
	// Load a user from the server.
	LoadUser(context.Context, LoadUserArg) (User, error)
	LoadUserByName(context.Context, LoadUserByNameArg) (User, error)
	// Load a user + device keys from the server.
	LoadUserPlusKeys(context.Context, LoadUserPlusKeysArg) (UserPlusKeys, error)
	// Load public keys for a user.
	LoadPublicKeys(context.Context, LoadPublicKeysArg) ([]PublicKey, error)
	// The list-tracking functions get verified data from the tracking statements
	// in the user's sigchain.
	//
	// If assertion is empty, it will use the current logged in user.
	ListTracking(context.Context, ListTrackingArg) ([]UserSummary, error)
	ListTrackingJSON(context.Context, ListTrackingJSONArg) (string, error)
	// Search for users who match a given query.
	Search(context.Context, SearchArg) ([]SearchResult, error)
	// Load all the user's public keys (even those in reset key families)
	// from the server with no verification
	LoadAllPublicKeysUnverified(context.Context, LoadAllPublicKeysUnverifiedArg) ([]PublicKey, error)
	ListTrackers2(context.Context, ListTrackers2Arg) (UserSummary2Set, error)
}

func UserProtocol(i UserInterface) rpc.Protocol {
	return rpc.Protocol{
		Name: "keybase.1.user",
		Methods: map[string]rpc.ServeHandlerDescription{
			"listTrackers": {
				MakeArg: func() interface{} {
					ret := make([]ListTrackersArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]ListTrackersArg)
					if !ok {
						err = rpc.NewTypeError((*[]ListTrackersArg)(nil), args)
						return
					}
					ret, err = i.ListTrackers(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"listTrackersByName": {
				MakeArg: func() interface{} {
					ret := make([]ListTrackersByNameArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]ListTrackersByNameArg)
					if !ok {
						err = rpc.NewTypeError((*[]ListTrackersByNameArg)(nil), args)
						return
					}
					ret, err = i.ListTrackersByName(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"listTrackersSelf": {
				MakeArg: func() interface{} {
					ret := make([]ListTrackersSelfArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]ListTrackersSelfArg)
					if !ok {
						err = rpc.NewTypeError((*[]ListTrackersSelfArg)(nil), args)
						return
					}
					ret, err = i.ListTrackersSelf(ctx, (*typedArgs)[0].SessionID)
					return
				},
				MethodType: rpc.MethodCall,
			},
			"loadUncheckedUserSummaries": {
				MakeArg: func() interface{} {
					ret := make([]LoadUncheckedUserSummariesArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]LoadUncheckedUserSummariesArg)
					if !ok {
						err = rpc.NewTypeError((*[]LoadUncheckedUserSummariesArg)(nil), args)
						return
					}
					ret, err = i.LoadUncheckedUserSummaries(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"loadUser": {
				MakeArg: func() interface{} {
					ret := make([]LoadUserArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]LoadUserArg)
					if !ok {
						err = rpc.NewTypeError((*[]LoadUserArg)(nil), args)
						return
					}
					ret, err = i.LoadUser(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"loadUserByName": {
				MakeArg: func() interface{} {
					ret := make([]LoadUserByNameArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]LoadUserByNameArg)
					if !ok {
						err = rpc.NewTypeError((*[]LoadUserByNameArg)(nil), args)
						return
					}
					ret, err = i.LoadUserByName(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"loadUserPlusKeys": {
				MakeArg: func() interface{} {
					ret := make([]LoadUserPlusKeysArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]LoadUserPlusKeysArg)
					if !ok {
						err = rpc.NewTypeError((*[]LoadUserPlusKeysArg)(nil), args)
						return
					}
					ret, err = i.LoadUserPlusKeys(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"loadPublicKeys": {
				MakeArg: func() interface{} {
					ret := make([]LoadPublicKeysArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]LoadPublicKeysArg)
					if !ok {
						err = rpc.NewTypeError((*[]LoadPublicKeysArg)(nil), args)
						return
					}
					ret, err = i.LoadPublicKeys(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"listTracking": {
				MakeArg: func() interface{} {
					ret := make([]ListTrackingArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]ListTrackingArg)
					if !ok {
						err = rpc.NewTypeError((*[]ListTrackingArg)(nil), args)
						return
					}
					ret, err = i.ListTracking(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"listTrackingJSON": {
				MakeArg: func() interface{} {
					ret := make([]ListTrackingJSONArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]ListTrackingJSONArg)
					if !ok {
						err = rpc.NewTypeError((*[]ListTrackingJSONArg)(nil), args)
						return
					}
					ret, err = i.ListTrackingJSON(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"search": {
				MakeArg: func() interface{} {
					ret := make([]SearchArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]SearchArg)
					if !ok {
						err = rpc.NewTypeError((*[]SearchArg)(nil), args)
						return
					}
					ret, err = i.Search(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"loadAllPublicKeysUnverified": {
				MakeArg: func() interface{} {
					ret := make([]LoadAllPublicKeysUnverifiedArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]LoadAllPublicKeysUnverifiedArg)
					if !ok {
						err = rpc.NewTypeError((*[]LoadAllPublicKeysUnverifiedArg)(nil), args)
						return
					}
					ret, err = i.LoadAllPublicKeysUnverified(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"listTrackers2": {
				MakeArg: func() interface{} {
					ret := make([]ListTrackers2Arg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]ListTrackers2Arg)
					if !ok {
						err = rpc.NewTypeError((*[]ListTrackers2Arg)(nil), args)
						return
					}
					ret, err = i.ListTrackers2(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
		},
	}
}

type UserClient struct {
	Cli rpc.GenericClient
}

func (c UserClient) ListTrackers(ctx context.Context, __arg ListTrackersArg) (res []Tracker, err error) {
	err = c.Cli.Call(ctx, "keybase.1.user.listTrackers", []interface{}{__arg}, &res)
	return
}

func (c UserClient) ListTrackersByName(ctx context.Context, __arg ListTrackersByNameArg) (res []Tracker, err error) {
	err = c.Cli.Call(ctx, "keybase.1.user.listTrackersByName", []interface{}{__arg}, &res)
	return
}

func (c UserClient) ListTrackersSelf(ctx context.Context, sessionID int) (res []Tracker, err error) {
	__arg := ListTrackersSelfArg{SessionID: sessionID}
	err = c.Cli.Call(ctx, "keybase.1.user.listTrackersSelf", []interface{}{__arg}, &res)
	return
}

// Load user summaries for the supplied uids.
// They are "unchecked" in that the client is not verifying the info from the server.
// If len(uids) > 500, the first 500 will be returned.
func (c UserClient) LoadUncheckedUserSummaries(ctx context.Context, __arg LoadUncheckedUserSummariesArg) (res []UserSummary, err error) {
	err = c.Cli.Call(ctx, "keybase.1.user.loadUncheckedUserSummaries", []interface{}{__arg}, &res)
	return
}

// Load a user from the server.
func (c UserClient) LoadUser(ctx context.Context, __arg LoadUserArg) (res User, err error) {
	err = c.Cli.Call(ctx, "keybase.1.user.loadUser", []interface{}{__arg}, &res)
	return
}

func (c UserClient) LoadUserByName(ctx context.Context, __arg LoadUserByNameArg) (res User, err error) {
	err = c.Cli.Call(ctx, "keybase.1.user.loadUserByName", []interface{}{__arg}, &res)
	return
}

// Load a user + device keys from the server.
func (c UserClient) LoadUserPlusKeys(ctx context.Context, __arg LoadUserPlusKeysArg) (res UserPlusKeys, err error) {
	err = c.Cli.Call(ctx, "keybase.1.user.loadUserPlusKeys", []interface{}{__arg}, &res)
	return
}

// Load public keys for a user.
func (c UserClient) LoadPublicKeys(ctx context.Context, __arg LoadPublicKeysArg) (res []PublicKey, err error) {
	err = c.Cli.Call(ctx, "keybase.1.user.loadPublicKeys", []interface{}{__arg}, &res)
	return
}

// The list-tracking functions get verified data from the tracking statements
// in the user's sigchain.
//
// If assertion is empty, it will use the current logged in user.
func (c UserClient) ListTracking(ctx context.Context, __arg ListTrackingArg) (res []UserSummary, err error) {
	err = c.Cli.Call(ctx, "keybase.1.user.listTracking", []interface{}{__arg}, &res)
	return
}

func (c UserClient) ListTrackingJSON(ctx context.Context, __arg ListTrackingJSONArg) (res string, err error) {
	err = c.Cli.Call(ctx, "keybase.1.user.listTrackingJSON", []interface{}{__arg}, &res)
	return
}

// Search for users who match a given query.
func (c UserClient) Search(ctx context.Context, __arg SearchArg) (res []SearchResult, err error) {
	err = c.Cli.Call(ctx, "keybase.1.user.search", []interface{}{__arg}, &res)
	return
}

// Load all the user's public keys (even those in reset key families)
// from the server with no verification
func (c UserClient) LoadAllPublicKeysUnverified(ctx context.Context, __arg LoadAllPublicKeysUnverifiedArg) (res []PublicKey, err error) {
	err = c.Cli.Call(ctx, "keybase.1.user.loadAllPublicKeysUnverified", []interface{}{__arg}, &res)
	return
}

func (c UserClient) ListTrackers2(ctx context.Context, __arg ListTrackers2Arg) (res UserSummary2Set, err error) {
	err = c.Cli.Call(ctx, "keybase.1.user.listTrackers2", []interface{}{__arg}, &res)
	return
}
