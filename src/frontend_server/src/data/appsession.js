
// TODO: remove this class - it's previosly used only for validation.
class AppSession
{
    constructor(request)
    {
        this.session = request.session;
    }

    get authorized()
    {
        return Boolean(this.session.authorized);
    }

    set authorized(value)
    {
        this.session.authorized = Boolean(value);
    }

    get username()
    {
        return '' + this.session.username;
    }

    set username(value)
    {
        this.session.username = '' + value;
    }
}

module.exports.AppSession = AppSession;
