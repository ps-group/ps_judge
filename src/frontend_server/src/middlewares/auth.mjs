export function checkAuth (req, res, next) 
{
    if (req.session.auth || req.path==='/login') 
    {
        next();
    } 
    else 
    {
       res.redirect("/login");
    }
}