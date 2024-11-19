const { evaluateUserCode } = require('./evaluator.js');
const utils = require('./ds_utils.js');

const userCode = `/**
 *
 * @param number n
 *
 * @returns number
 */
function factorial(n) {
    // TODO: Implement this function
}`;

const testCases = [{"parameters":["0"],"expected_output":"1"},{"parameters":["3"],"expected_output":"6"},{"parameters":["7"],"expected_output":"5040"}];
const functionName = "factorial";

const results = evaluateUserCode(userCode, testCases, functionName);
console.log(JSON.stringify(results, null, 2));
