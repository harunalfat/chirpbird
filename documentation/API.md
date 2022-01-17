## API Documentation

For Chirpbird app, we use [Centrifuge](https://github.com/centrifugal/centrifuge) as backbone of our pub/sub service. Your client implementation may available within [Centrifugal repositories](https://github.com/centrifugal). If not, sadly, you need to write your own client, or change to existing stack that are supported :P.

### Important

This app is created for the purpose of basic requirements. Make sure it can send message between users, privately and publicly within a group. Right now, the authentication relies on the userId only. THIS IS NOT a proper way!

BASE_URL = `https://harunalfat.site`

### User Registration
Send GET request to BASE_URL with payload 
```
{
    "username": "Alfat"
}
```

Success response example:
```
{
    "data": {
        "id": "0e16b1a1-36ad-4532-a5a5-33c58a2e866b",
        "createdAt": "2022-01-15T02:15:46.484560878Z",
        "updatedAt": "0001-01-01T00:00:00Z",
        "username": "Alfat",
        "channels": [
            {
                "id": "xxxxx-yyyyy-zzzzz",
                "name": "Lobby"
            }
        ]
    },
    "errors": null
}
```
Failed response example:
```
{
    "data": null,
    "errors": ["error message 1", "error message 2"]
}
```

## Establish connection
You only need to provide our BASE_URL with url query `userId` with your user ID above, within the client constructor.
```
// by using centrifuge.js
instance = new Centrifuge(`${BASE_URL}?userId=${userId}`);
instance.connect();
instance.on("connect", func(ctx) {
    console.log("I'm connected")
    doSomethingWithMessage(ctx)
});
```
Please don't use protobuf protocol  (default one using JSON), as it currently not supported on Chirpbird.

Centrifuge object is an EventEmitter. You can listen to `connect` or `disconnect`, to do action regarding your connection status.

## Joining Channel and Listen
When connected, you can listen incoming messages from listed channels that is responded by `User Registration` part above.

Simply loop through the user's channels to subscribe to all

```
for (const channel of user.channels) {
    const subscription = instance.subscribe(channel.id)
    subscription.on("publish", func(ctx) {
        console.log("Do something")
        doSomethingWithMessage(ctx)
    })
}
```

If you want to fetch the channel previous messages, use RPC method name `fetch_messages`
```
const promiseResult = instance.namedRPC("fetch_messages", {
    data: "the-channel-id"
})

const result = await promiseResult
```

`result` will be structured as
```
{
    data: [
        {
            id: "xxxxx",
            createdAt: "2022-01-15T05:48:58.491Z",
            data: "Testing the message",
            sender: {
                id: "sender-id-here",
                username: "my-username"
            }
        }
    ]
}
```
## Publish Message to Channel
Publishing a message will store it to server database and also publishing it to other subscriber of that channel
```
await instance.publish("our-lobby-id", {
    data: "I try to send you a message"
})
```

at the receiving end,
```
const subscription = instance.subscribe("our-lobby-id")
subscription.on("publish", func(ctx) {
    console.log("Do something")
    console.log(ctx)
    doSomethingWithMessage(ctx)
})
```
You can do something within the callback, maybe add it to your DOM object, push the message list, or any kind of thing

The `ctx` object is structured as
```
{
    data: {
        data: "I try to send you a message",
        sender: {
            id: "the-sender-id",
            username: "the-sender-username"
        }
        channel: {
            id: "the-channel-id",
            name: "Lobby"
        }
    }
}
```

## Add New Channel
If you want your user to participate to more channels, use RPC method name `create_channel`
```
const channel = {
    name: "Heal the world",
    isPrivate: false,
    participants: [], // no need to add yourself as participant
}

// or for private channel
const channel = {
    name: "Heal the world",
    isPrivate: false,
    participants: [
        {
            id: "my-friend-id",
            username: "My Friend User Name"
        }
    ],
}

const promiseResult = instance.namedRPC("create_channel", {
    data: channel
})

const result = await promiseResult
```
`result` will be structured as
```
{
    data: {
        id: "given-id",
        name: "Heal the world",
        isPrivate: false,
        createdAt: "2022-01-15T05:48:58.491Z",
        creatorId: "one-who-create-uuid"
    }
}
```
It will create the channel if it's not exist yet. If the name already exist, it simply add you as a participant of that channel, and can listen event to that.

Either subscribe to or create existing channel name that is set to private will result in error. Except you already the participant

## Search Users
Maybe you want to get some user's ID to add them into private channel? You can use RPC method name `search_users`
```
const promiseResult = instance.namedRPC("search_users", {
    data: "Chir",
})

const result = await promiseResult
```
`result` will be structured as
```
{
    data: [
        {
            id: "given-id",
            username: "Chirpbird"
        },
        {
            id: "given-id2",
            username: "chiri"
        }
    ]
}
```

## Listen When Added to New Channel
When you are being added to new channel, the publication came from the server itself, so you will listen to server notification
```
instance.on("message", func(ctx) { // <--- Please note, it listen from "instance", not "subscription"
    console.log(ctx)
    if (ctx.channel === "NEW_CHANNEL") { // <--- Please note, ctx.channel here
        // is the naming schema of Centrifuge to let client know
        // what context the server message are,
        // it is not related to the "group chat" channel

        const newAnnouncedChannel = ctx.data;
        const subscription = instance.subscribe(newAnnouncedChannel.id)
        subscription.on("publish", , func(ev) {
            doSomethingWithMessage(ev);
        }
    })
    }
})
```
