const amqplib = require('amqplib');
const exchangeBuildFinished = 'psjudge-build-finished'
const queueBuildFinishedFrontend = 'psjudge-build-finished-frontend'

class MessageRouter
{
    constructor()
    {
        this._launchPromise = this.launch();
        this._connection = null;
    }

    async launch()
    {
        this._connection = await amqplib.connect();
        this._channel = await this._connection.createChannel();
        this._channel.assertQueue(queueBuildFinishedFrontend, { durable: true });
        this._channel.assertExchange(exchangeBuildFinished, 'fanout', { durable: false });
        this._channel.bindQueue(queueBuildFinishedFrontend, exchangeBuildFinished, '');
    }

    async close() {
        await this._launchPromise;
        if (this._connection)
        {
            this._connection.close();
        }
    }

    async consumeBuildFinished(cb)
    {
        await this._consumeQueue(queueBuildFinishedFrontend, cb);
    }

    async _consumeQueue(queueName, cb)
    {
        await this._launchPromise;
        const wrappedCb = (msg) => {
            try
            {
                const json = msg.content.toString();
                const value = JSON.parse(json);
                cb(value);
            }
            catch (error)
            {
                console.error('exception while processing RabbitMQ message', error);
            }
        };
        const options = {
            noAck: true
        };
        this._channel.consume(queueName, wrappedCb, options);
    }

    async _publish(exchange, value)
    {
        await this._launchPromise;
        const json = JSON.stringify(value);
        this._channel.publish(exchange, '', new Buffer(json));
    }
}

module.exports.MessageRouter = MessageRouter;
