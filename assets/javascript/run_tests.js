const { evaluateUserCode } = require('./evaluator.js');
const userCode = `/**
 *
 * @param string s
 *
 * @returns boolean
 */
function isPalindrome(s) {
     if (!tree1 && !tree2) return true; // Both trees are null
    if (!tree1 || !tree2) return false; // One tree is null, the other isn't
    if (tree1.val !== tree2.val) return false; // Node values differ

    // Recursively compare left and right subtrees
    return compareTrees(tree1.left, tree2.left) && compareTrees(tree1.right, tree2.right);
}















`;

const testCases = [{"parameters":["'madam'"],"expected_output":"true"},{"parameters":["'hello'"],"expected_output":"false"},{"parameters":["'A man a plan a canal Panama'"],"expected_output":"true"}];
const functionName = "isPalindrome";

const results = evaluateUserCode(userCode, testCases, functionName);
console.log(JSON.stringify(results, null, 2));
