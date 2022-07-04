<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [gaia/usc/v1beta1/usc.proto](#gaia/usc/v1beta1/usc.proto)
    - [Params](#gaia.usc.v1beta1.Params)
    - [RedeemEntry](#gaia.usc.v1beta1.RedeemEntry)
    - [RedeemEntryOperation](#gaia.usc.v1beta1.RedeemEntryOperation)
    - [RedeemingQueueData](#gaia.usc.v1beta1.RedeemingQueueData)
    - [TokenMeta](#gaia.usc.v1beta1.TokenMeta)
  
- [gaia/usc/v1beta1/genesis.proto](#gaia/usc/v1beta1/genesis.proto)
    - [GenesisState](#gaia.usc.v1beta1.GenesisState)
  
- [gaia/usc/v1beta1/query.proto](#gaia/usc/v1beta1/query.proto)
    - [QueryParamsRequest](#gaia.usc.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#gaia.usc.v1beta1.QueryParamsResponse)
    - [QueryPoolRequest](#gaia.usc.v1beta1.QueryPoolRequest)
    - [QueryPoolResponse](#gaia.usc.v1beta1.QueryPoolResponse)
    - [QueryRedeemEntryRequest](#gaia.usc.v1beta1.QueryRedeemEntryRequest)
    - [QueryRedeemEntryResponse](#gaia.usc.v1beta1.QueryRedeemEntryResponse)
  
    - [Query](#gaia.usc.v1beta1.Query)
  
- [gaia/usc/v1beta1/tx.proto](#gaia/usc/v1beta1/tx.proto)
    - [MsgMintUSC](#gaia.usc.v1beta1.MsgMintUSC)
    - [MsgMintUSCResponse](#gaia.usc.v1beta1.MsgMintUSCResponse)
    - [MsgRedeemCollateral](#gaia.usc.v1beta1.MsgRedeemCollateral)
    - [MsgRedeemCollateralResponse](#gaia.usc.v1beta1.MsgRedeemCollateralResponse)
  
    - [Msg](#gaia.usc.v1beta1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="gaia/usc/v1beta1/usc.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## gaia/usc/v1beta1/usc.proto



<a name="gaia.usc.v1beta1.Params"></a>

### Params
Params defines the parameters for the x/usc module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `redeem_dur` | [google.protobuf.Duration](#google.protobuf.Duration) |  | redeem_dur defines USC -> collateral redeem duration (how long does it takes to convert). |
| `max_redeem_entries` | [uint32](#uint32) |  | max_redeem_entries is a max number of concurrent redeem operations per account. |
| `collateral_metas` | [TokenMeta](#gaia.usc.v1beta1.TokenMeta) | repeated | collateral_metas defines a set of collateral token metas that are supported by the module. |






<a name="gaia.usc.v1beta1.RedeemEntry"></a>

### RedeemEntry
RedeemEntry defines a redeeming queue object entry.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is a redeem target account. |
| `operations` | [RedeemEntryOperation](#gaia.usc.v1beta1.RedeemEntryOperation) | repeated | operations are redeem operations that are active. |






<a name="gaia.usc.v1beta1.RedeemEntryOperation"></a>

### RedeemEntryOperation
RedeemEntryOperation defines a single redeeming queue operation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `creation_height` | [int64](#int64) |  | creation_height is the height which the redeeming took place. |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | completion_time is the unix time for redeeming completion. |
| `collateral_amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | collateral_amount are collateral tokens to redeem. |






<a name="gaia.usc.v1beta1.RedeemingQueueData"></a>

### RedeemingQueueData
RedeemingQueueData defines the redeeming queue value object (completionTime timestamp is used as a key for the queue).
Object is used to link queue data with a corresponding RedeemEntry object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `addresses` | [string](#string) | repeated |  |






<a name="gaia.usc.v1beta1.TokenMeta"></a>

### TokenMeta
TokenMeta defines USC / collateral token metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | denom is the sdk.Coin denomination (ibc / native tokens). |
| `decimals` | [uint32](#uint32) |  | decimals is a number of decimals for the sdk.Coin amount value. Example: coin amount of 1.0usdt with 3 decimals -> 1000. |
| `description` | [string](#string) |  | description is an optional token description (IBC source info for example). |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="gaia/usc/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## gaia/usc/v1beta1/genesis.proto



<a name="gaia.usc.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the x/usc module genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#gaia.usc.v1beta1.Params) |  | params are all the module parameters. |
| `redeem_entries` | [RedeemEntry](#gaia.usc.v1beta1.RedeemEntry) | repeated | redeem_entries are active redeeming queue objects. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="gaia/usc/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## gaia/usc/v1beta1/query.proto



<a name="gaia.usc.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is request type for the Query/Params RPC method.






<a name="gaia.usc.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#gaia.usc.v1beta1.Params) |  | params are the current module parameters. |






<a name="gaia.usc.v1beta1.QueryPoolRequest"></a>

### QueryPoolRequest
QueryPoolRequest is request type for Query/Pool RPC method.






<a name="gaia.usc.v1beta1.QueryPoolResponse"></a>

### QueryPoolResponse
QueryPoolResponse is response type for the Query/Pool RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `active_pool` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | active_pool is the Active pool supply that could be used for redeeming. |
| `redeeming_pool` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | redeeming_pool is the Redeeming pool supply that is currently being redeemed. |






<a name="gaia.usc.v1beta1.QueryRedeemEntryRequest"></a>

### QueryRedeemEntryRequest
QueryRedeemEntryRequest is request type for the Query/RedeemEntry RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the target account address. |






<a name="gaia.usc.v1beta1.QueryRedeemEntryResponse"></a>

### QueryRedeemEntryResponse
QueryRedeemEntryResponse is response type for the Query/RedeemEntry RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [RedeemEntry](#gaia.usc.v1beta1.RedeemEntry) |  | entry is the redeeming entry for an account. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="gaia.usc.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Pool` | [QueryPoolRequest](#gaia.usc.v1beta1.QueryPoolRequest) | [QueryPoolResponse](#gaia.usc.v1beta1.QueryPoolResponse) | Pool queries the collateral balance pool info. | GET|/gaia/usc/v1beta1/pool|
| `Params` | [QueryParamsRequest](#gaia.usc.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#gaia.usc.v1beta1.QueryParamsResponse) | Params queries the module parameters. | GET|/gaia/usc/v1beta1/params|
| `RedeemEntry` | [QueryRedeemEntryRequest](#gaia.usc.v1beta1.QueryRedeemEntryRequest) | [QueryRedeemEntryResponse](#gaia.usc.v1beta1.QueryRedeemEntryResponse) | RedeemEntry queries a redeem entry for an account. | GET|/gaia/usc/v1beta1/redeem-entry|

 <!-- end services -->



<a name="gaia/usc/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## gaia/usc/v1beta1/tx.proto



<a name="gaia.usc.v1beta1.MsgMintUSC"></a>

### MsgMintUSC
MsgMintUSC defines a SDK message for the Msg/MintUSC request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the Bech32-encoded address of the target account. |
| `collateral_amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | collateral_amount are collateral token that should be exchanged to USC. |






<a name="gaia.usc.v1beta1.MsgMintUSCResponse"></a>

### MsgMintUSCResponse
MsgMintUSCResponse defines the Msg/MintUSC response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `minted_amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | minted_amount is a minted USC token. |
| `collaterals_amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | collaterals_amount are collateral tokens used for the mint. |






<a name="gaia.usc.v1beta1.MsgRedeemCollateral"></a>

### MsgRedeemCollateral
MsgRedeemCollateral defines a SDK message for the Msg/RedeemCollateral request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the Bech32-encoded address of the target account. |
| `usc_amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | usc_amount is the USC token that should be exchanged to collateral tokens. |






<a name="gaia.usc.v1beta1.MsgRedeemCollateralResponse"></a>

### MsgRedeemCollateralResponse
MsgMintUSCResponse defines the Msg/RedeemCollateral response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `burned_amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | burned_amount is the USC token converted amount (might be LT the requested amount). |
| `redeemed_amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | redeemed_amount are collateral tokens that are transferred to an account after the redeeming timout. |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | completion_time defines the redeeming period end time. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="gaia.usc.v1beta1.Msg"></a>

### Msg
Msg defines the x/usc Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `MintUSC` | [MsgMintUSC](#gaia.usc.v1beta1.MsgMintUSC) | [MsgMintUSCResponse](#gaia.usc.v1beta1.MsgMintUSCResponse) | MintUSC defines a method for converting collateral coins into USC coin. | |
| `RedeemCollateral` | [MsgRedeemCollateral](#gaia.usc.v1beta1.MsgRedeemCollateral) | [MsgRedeemCollateralResponse](#gaia.usc.v1beta1.MsgRedeemCollateralResponse) | RedeemCollateral defines a method for converting USC coin into collateral coins from the module pool. | |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

