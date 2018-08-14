/**
 * @param {any} value
 * @return {string} type of the value, 'null' for null, or class name if value is objects
 */
function getTypeOf(value)
{
    const type = typeof value;
    if (type == 'object')
    {
        if (value === null)
        {
            return 'null';
        }
        return Object.prototype.toString.call(value).match(/^\[object (.*)\]$/)[1];
    }
    return type;
}

/**
 * @param value - any value that is valid only if it's integer.
 * @returns {number}
 */
export function verifyInt(value)
{
    const type = getTypeOf(value);
    if (type == 'Number')
    {
        value = value.valueOf();
    }
    else if (type != 'number')
    {
        throw new Error('value must be integer but has type ' + type);
    }
    if (!Number.isInteger(value))
    {
        throw new Error('value must be integer but is floating-point');
    }
    return value;
}

/**
 * @param value - any value that is valid only if it's string
 */
export function verifyString(value)
{
    const type = getTypeOf(value);
    if (type == 'String')
    {
        value = value.valueOf();
    }
    else if (type != 'string')
    {
        throw new Error('value must be string but has type ' + type);
    }
    return value;
}

/**
 * @param value - any value that is valid only if it's array
 */
export function verifyArray(value)
{
    const type = getTypeOf(value);
    if (type == 'Array')
    {
        return value;
    }
    throw new Error('value must be Array but has type ' + type);
}
