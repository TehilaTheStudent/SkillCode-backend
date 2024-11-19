const { evaluateUserCode } = require('./evaluator.js');


const userCode = `/**
 *
 * @param Array<number> arr
 *
 * @param number target
 *
 * @returns number
 */
function binarySearch(arr, target) {
    // TODO: Implement this function
}`;

const testCases = [{"parameters":["[1, 2, 3, 4, 5]","3"],"expected_output":"2"},{"parameters":["[1, 2, 3, 4, 5]","6"],"expected_output":"-1"}];
const functionName = "binarySearch";

const results = evaluateUserCode(userCode, testCases, functionName);
console.log(JSON.stringify(results, null, 2));
