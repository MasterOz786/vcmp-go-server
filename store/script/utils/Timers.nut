Timer <- {
	Timers = {}

	function Create(environment, listener, interval, repeat, ...)
	{
		vargv.insert(0, environment);

		local TimerInfo = {
			Environment = environment,
			Listener = listener,
			Interval = interval,
			Repeat = repeat,
			Args = vargv,
			LastCall = Script.GetTicks(),
			CallCount = 0
		};

		local hash = split(TimerInfo.tostring(), ":")[1].slice(3, -1).tointeger(16);
		hash = hash + TimerInfo.Args[1].tostring();
		Timers.rawset(hash, TimerInfo);
		return hash;
	}

	function Destroy(hash)
	{
		if (Timers.rawin(hash)) {
			Timers.rawdelete(hash);
		}
	}

	function Process()
	{
		local CurrTime = Script.GetTicks();
		foreach (hash, tm in Timers)
		{
			if (tm != null && CurrTime - tm.LastCall >= tm.Interval)
			{
				tm.CallCount++;
				tm.LastCall = CurrTime;
				tm.Listener.pacall(tm.Args);
				if (tm.Repeat != 0 && tm.CallCount >= tm.Repeat) {
					Timers.rawdelete(hash);
				}
			}
		}
	}
};
