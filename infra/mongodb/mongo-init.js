db.createUser(
	{
		user: "mongotest",
		pwd: "mongotest",
		roles: [
			{
				role: "readWrite",
				db: "mongotest"
			}
		]
	}
);

