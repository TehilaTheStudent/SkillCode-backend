const dsUtils = require('./ds_utils');

function listyToType(stringyListyRep, abstractType) {
    /**
     * Convert a stringy listy representation into the appropriate data structure
     * based on the given abstract type.
     *
     * Args:
     *     stringyListyRep (str): The stringy representation of the data.
     *     abstractType (dict): The abstract type definition with "type" and optional "typeChildren".
     *
     * Returns:
     *     Any: The converted data structure.
     */
    let listyRep;
    try {
        // Safely parse the stringy representation into JavaScript structures
        listyRep = JSON.parse(stringyListyRep);
    } catch (e) {
        throw new Error(`Failed to parse input: ${stringyListyRep}. Error: ${e.message}`);
    }

    const baseType = abstractType.type;
    const typeChildren = abstractType.typeChildren;

    // Handle atomic types
    if (["Integer", "Double", "String", "Boolean"].includes(baseType)) {
        // Directly return the value for atomic types
        return listyRep;
    }

    // Handle composite types
    if (baseType === "Array") {
        // Recursively convert each element using the child type
        return listyRep.map(item => listyToType(JSON.stringify(item), typeChildren));
    }

    if (baseType === "Matrix") {
        // Recursively convert each row using the child type
        return listyRep.map(row => listyToType(JSON.stringify(row), typeChildren));
    }

    if (baseType === "TreeNode") {
        // Use the utility function to generate a TreeNode
        let tree = dsUtils.generateTree(listyRep);
        if (typeChildren) {
            // Convert the value and children using the child type
            tree.val = listyToType(JSON.stringify(tree.val), typeChildren);
            if (tree.left) {
                tree.left = listyToType(JSON.stringify(dsUtils.exportTree(tree.left)), abstractType);
            }
            if (tree.right) {
                tree.right = listyToType(JSON.stringify(dsUtils.exportTree(tree.right)), abstractType);
            }
        }
        return tree;
    }

    if (baseType === "ListNode") {
        // Use the utility function to generate a ListNode
        let linkedList = dsUtils.generateLinkedList(listyRep);
        if (typeChildren) {
            // Convert each node value using the child type
            let current = linkedList;
            while (current) {
                current.val = listyToType(JSON.stringify(current.val), typeChildren);
                current = current.next;
            }
        }
        return linkedList;
    }

    if (baseType === "Graph") {
        // Use the utility function to generate a Graph
        let graph = dsUtils.generateGraph(listyRep);
        if (typeChildren) {
            // Convert each node and its neighbors using the child type
            let newAdjList = {};
            graph.forEach((node, u) => {
                let convertedU = listyToType(JSON.stringify(u), typeChildren);
                newAdjList[convertedU] = node.neighbors.map(neighbor => listyToType(JSON.stringify(neighbor.val), typeChildren));
            });
            graph.adjList = newAdjList;
        }
        return graph;
    }

    throw new Error(`Unsupported type: ${baseType}`);
}

module.exports = {
    listyToType,
};
