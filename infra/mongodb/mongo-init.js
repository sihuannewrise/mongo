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
// db.getSiblingDB("myDatabase")
// 	.createUser({
// 		user: "myUser",
// 		pwd: "myPassword",
// 		roles: [{ role: "readWrite", db: "myDatabase" }],
// 		passwordDigestor: "server"
// 	});
