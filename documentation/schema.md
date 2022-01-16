## Chirpbird Schema
---

This is a relation for each entity that are used within the app.
### Base
Described as
```
type Base struct {
	ID        string    `json:"id" bson:"id"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
```

It simply provide us with base fields that will be need by any other entities.

### User
We describe it as
```
type User struct {
	Base     `bson:",inline"`
	Username string   `json:"username,omitempty" bson:"username,omitempty"`
	Channels Channels `json:"channels,omitempty" bson:"channels,omitempty"`
}
```
For the sake of simplicity, we only use username. `Channels` is useful for tracking what channel the user participates in. When user add more channel, this field also will be appended with the new one.

### Channel
Described as
```
type Channel struct {
	Base           `bson:",inline"`
	Name           string `json:"name,omitempty" bson:"name,omitempty"`
	CreatorID      string `json:"creatorId,omitempty" bson:"creatorId,omitempty"`
	IsPrivate      bool   `json:"isPrivate" bson:"isPrivate"`
	Participants   []User `json:"participants,omitempty" bson:"participants,omitempty"`
	HashIdentifier string `json:"hashIdentifier,omitempty" bson:"hashIdentifier,omitempty"`
}
```
`CreatorID` is using string, for no particular reason. I guess it'll be better if we are using `User` directly as `Creator` field.

`IsPrivate` will flag the channel as private, no additional member can join else than one original users registered when it was created. This app currently not supporting invitation for more than 2 users in private channel group.

`Participants` as the reversed role of `Channels` in `User`. Will list all of its participated users.

`HashIdentifier` is my simple hack to support unique index with the channel name. In example, `Hello World` channel name already exist, and someone create `HELLO WORLD`, it will give same `HashIdentifier`, so it'll return former one. Mongodb supports unique with case insensitive by using `collation`, but I rather work with this for such a small case like this app than debugging why it doesn't work later on.

### Message
Described as
```
type Message struct {
	Base      `bson:",inline"`
	Sender    User        `json:"sender,omitempty" bson:"sender,omitempty"`
	ChannelID string      `json:"channelId,omitempty" bson:"channelId,omitempty"`
	Data      interface{} `json:"data,omitempty" bson:"data,omitempty"`
}
```
`Sender` is the user who send the message.
`ChannelID` I guess will be more fitted to directly link it as `Channel`. The `Data` originally can be more than string, but yeah, for simplicity.

## Relationship Between Entity

Different than relational database (RDB), there is no such things as relationship between "tables". To "relate" entities, it used collection aggregation. So in my Mongodb, the strategy that I use is:

- Save list of `User` in its own collection called `users`
- Save list of `Channel` in its own collection called `channels`
- Save list of `Message` appropriate to its own channel with its own collection, called `channel:xxxxx-aaaaa-bbbb-yyyy`

How do I store the info? I store them all :P. I store them as it is.
One example of `User` in `users`
```
{
    "_id" : ObjectId("61e1b0446bd9377f2fe5e4f2"),
    "id" : "e409699e-fc75-4a07-acb7-bece303820d5",
    "createdAt" : ISODate("2022-01-14T17:17:56.864Z"),
    "username" : "Alfat",
    "channels" : [ 
        {
            "id" : "f217d2fb-235b-4b63-ae73-6466eebbb6c9",
            "name" : "Lobby"
        }, 
        {
            "id" : "beff5424-978e-43d2-9616-04e582a8084a",
            "name" : "Alfat & Sopo"
        }, 
        {
            "id" : "18887f14-334d-4df9-820f-5bbf42d1c793",
            "name" : "Alfat & Sopoja"
        }, 
        {
            "id" : "211ad980-d09f-4e64-a030-a947b53ac1e1",
            "name" : "Alfat & Jekarda"
        }
    ]
}
```
One example of `Channel` in `channels`
```
{
    "_id" : ObjectId("61e273421ab5c70015c93d69"),
    "id" : "18887f14-334d-4df9-820f-5bbf42d1c793",
    "name" : "Alfat & Sopoja",
    "creatorId" : "86256088-3f4b-4f86-baf9-7712ba746f1f",
    "isPrivate" : true,
    "participants" : [ 
        {
            "id" : "86256088-3f4b-4f86-baf9-7712ba746f1f",
            "username" : "Sopoja"
        }, 
        {
            "id" : "e409699e-fc75-4a07-acb7-bece303820d5",
            "username" : "Alfat"
        }
    ],
    "hashIdentifier" : "713fd678b9df81c2a2014c15e50a41e312e18af871243db4d5feb85caf8d24aa"
}
```
One example of `Message` in `messages`
```
{
    "_id" : ObjectId("61e29d26ffc7a9d45c0d0eaf"),
    "id" : "46dd2914-b9f5-4b48-b0d8-99e6952407bd",
    "createdAt" : ISODate("2022-01-15T10:08:38.590Z"),
    "sender" : {
        "id" : "985bc94b-ea2c-45f1-909b-27ca0375be13",
        "username" : "Jaka"
    },
    "channelId" : "f217d2fb-235b-4b63-ae73-6466eebbb6c9",
    "data" : "woit"
}
```
You will notice that are many duplication between each collection, like `User` name reappeared in `Channel` as participant and `Message` as sender. This is a drawback and also the feature of key value store like Mongodb. In RDB, this will be slammed as bad practice, and I wholeheartedly aggree, each table should be in normalized condition in RDB. With using key value, this is a different case. The goal of it, space aren't the issue, querying is, so with Mongodb, we trade complexity of querying something with, just store them all in one place. In the app, my example is the chat history, it uses its own database for each channel destination. Someone can even make something like collection key `someone-user-id:channels` and store each of its own subscribed channel to it. Or maybe `some-channel-id:participants` and adding a `Channel` participants to it.

You can still manage what fields need to be stored, in my case, I only store `Channel` ID and Name on `users`' document, no need for other things. But beware of data consistencies, in this app case, if user changing his username, the message history sender won't be altered, except we want to, because the data is not related for each `Collection`.