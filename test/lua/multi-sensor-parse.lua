Actions = {
	function(args)
		local dataT, err = json:J2T(args)
		if (err ~= nil) then
			Throw(err)
			return true, args
		end
		local params = {}
		for _, value in pairs(dataT) do
			params[value['tag']] = value.value
		end
		local json = json:T2J({
			id = time:TimeMs(),
			method = "thing.event.property.post",
			params = params,
		})
		Debug(json)
		return true, args
	end
}
