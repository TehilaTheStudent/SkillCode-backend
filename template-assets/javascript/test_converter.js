const assert = require('assert');
const dsUtils = require('./ds_utils');
const { TreeNode, Graph, ListNode } = dsUtils;
const { listyToType } = require('./converter');

describe('listyToType Integration Tests', function() {
    it('should convert to TreeNode', function() {
        const stringyInput = "[1, 2, 3, 1, 1, 4, 5]";
        const abstractType = {
            type: "TreeNode",
            typeChildren: {
                type: "Integer"
            }
        };

        const tree = listyToType(stringyInput, abstractType);
        assert(tree instanceof TreeNode);
        assert.deepStrictEqual(dsUtils.exportTree(tree), [1, 2, 3, 1, 1, 4, 5]);
    });

    it('should convert to Graph', function() {
        const stringyInput = "[[1, 2], [2, 3], [3, 1]]";
        const abstractType = {
            type: "Graph",
            typeChildren: {
                type: "Integer"
            }
        };

        const graph = listyToType(stringyInput, abstractType);
        console.log(graph);
        console.log(dsUtils.exportGraph(graph));
        assert(graph instanceof Graph);
        assert.deepStrictEqual(dsUtils.exportGraph(graph), [[1, 2], [2, 3], [3, 1]]);
    });

    it('should convert to ListNode', function() {
        const stringyInput = "[1, 2, 3, 4]";
        const abstractType = {
            type: "ListNode",
            typeChildren: {
                type: "Integer"
            }
        };

        const linkedList = listyToType(stringyInput, abstractType);
        assert(linkedList instanceof ListNode);
        assert.deepStrictEqual(dsUtils.exportLinkedList(linkedList), [1, 2, 3, 4]);
    });

    it('should convert to Array', function() {
        const stringyInput = "[1, 2, 3, 4]";
        const abstractType = {
            type: "Array",
            typeChildren: {
                type: "Integer"
            }
        };

        const array = listyToType(stringyInput, abstractType);
        assert(Array.isArray(array));
        assert.deepStrictEqual(array, [1, 2, 3, 4]);
    });

    it('should convert to Matrix', function() {
        const stringyInput = "[[1, 2], [3, 4], [5, 6]]";
        const abstractType = {
            type: "Matrix",
            typeChildren: {
                type: "Integer"
            }
        };

        const matrix = listyToType(stringyInput, abstractType);
        assert(Array.isArray(matrix));
        assert.deepStrictEqual(matrix, [[1, 2], [3, 4], [5, 6]]);
    });

    it('should throw error for invalid type', function() {
        const stringyInput = "[1, 2, 3]";
        const abstractType = { type: "Unknown" };
        assert.throws(() => listyToType(stringyInput, abstractType), Error);
    });
});
