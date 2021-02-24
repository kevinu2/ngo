package constant

var (
	CodeUnLogin = 2013

	CodeRspTradeClosed = 4001 //休市

	CodeRspTradeMaintain = 4002 //维护

	CodeRspSystemError = 5000 //系统异常

	CodeRspParamQuantity = 6012 //订单信息错误

	CodeRspTradeCoinSymbolError = 6001 //交易交易币种错误

	CodeRspTradeAmountError = 6002 //交易数量错误

	CodeRspTradeConfigError = 6003 //交易配置错误

	CodeRspTradeDataError = 6004 //交易数据非法

	CodeRspTradeInitDataError = 6005 //初始化交易数据异常

	CodeRspTradeOverBalanceError = 6006 //交易超出余额

	CodeRspTradeOverMaxLimitError = 6007 //交易超出系统最大限定

	CodeRspTradeOverMinLimitError = 6008 //交易超出系统最小限定

	CodeRspOrderLimitError = 6009 //订单超出最大

	CodeRspParamProductError = 6011 //产品信息错误

	CodeRspParamOrderError = 6012 //订单信息错误

	CodeRspParamParseError = 6021 //参数解析错误

	CodeRspNotRequiredDepositError = 6022 //无需充值

	CodeRspNotEnoughRewardsError = 6023 //赠金不足

	CodeRspNotAllowedCoinIdError = 6024 //币种不允许

	CodeRspInternalError = 6000 //内部错误

	CodeRspTradeOverBuyProtect = 7000 //今日因买入额过大已触发流动性保护机制，您的买单已进入列队执行序列等待成交

	CodeRspTradeOverTotalBalanceError = 7001 //委托的累计usdt金额大于当前可以usdt余

)

const (
	NoErr         uint16 = 0
	SystemCodeErr uint16 = 5000
)
