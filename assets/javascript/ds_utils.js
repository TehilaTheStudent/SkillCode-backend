class GraphNode {
    /**
     * A node in a graph.
     * @param {any} val - The value of the node.
     * @param {GraphNode[]} neighbors - List of neighbor nodes.
     */
    constructor(val = null, neighbors = []) {
        this.val = val;
        this.neighbors = neighbors;
    }

    toString() {
        return `GraphNode(val=${this.val})`;
    }
}

class TreeNode {
    /**
     * A node in a binary tree.
     * @param {any} val - The value of the node.
     * @param {TreeNode|null} left - Left child.
     * @param {TreeNode|null} right - Right child.
     */
    constructor(val = null, left = null, right = null) {
        this.val = val;
        this.left = left;
        this.right = right;
    }

    toString() {
        return `TreeNode(val=${this.val})`;
    }
}

class ListNode {
    /**
     * A node in a singly linked list.
     * @param {any} val - The value of the node.
     * @param {ListNode|null} next - The next node.
     */
    constructor(val = null, next = null) {
        this.val = val;
        this.next = next;
    }

    toString() {
        return `ListNode(val=${this.val})`;
    }
}

// Generate a binary tree from a list of values
function generateTree(values) {
    if (!values || values.length === 0) return null;

    const root = new TreeNode(values[0]);
    const queue = [root];
    let i = 1;

    while (queue.length > 0 && i < values.length) {
        const current = queue.shift();

        if (values[i] !== null) {
            current.left = new TreeNode(values[i]);
            queue.push(current.left);
        }
        i++;

        if (i < values.length && values[i] !== null) {
            current.right = new TreeNode(values[i]);
            queue.push(current.right);
        }
        i++;
    }

    return root;
}

// Export a binary tree to a list of values
function exportTree(root) {
    if (!root) return [];

    const result = [];
    const queue = [root];

    while (queue.length > 0) {
        const current = queue.shift();
        if (current) {
            result.push(current.val);
            queue.push(current.left);
            queue.push(current.right);
        } else {
            result.push(null);
        }
    }

    // Remove trailing null values
    while (result[result.length - 1] === null) {
        result.pop();
    }

    return result;
}

// Generate a singly linked list from a list of values
function generateLinkedList(values) {
    if (!values || values.length === 0) return null;

    const head = new ListNode(values[0]);
    let current = head;

    for (let i = 1; i < values.length; i++) {
        current.next = new ListNode(values[i]);
        current = current.next;
    }

    return head;
}

// Export a singly linked list to a list of values
function exportLinkedList(head) {
    const result = [];
    let current = head;

    while (current) {
        result.push(current.val);
        current = current.next;
    }

    return result;
}

// Generate a graph from a list of edges
function generateGraph(edges, directed = false) {
    const nodes = new Map();

    edges.forEach(([u, v]) => {
        if (!nodes.has(u)) nodes.set(u, new GraphNode(u));
        if (!nodes.has(v)) nodes.set(v, new GraphNode(v));

        nodes.get(u).neighbors.push(nodes.get(v));
        if (!directed) {
            nodes.get(v).neighbors.push(nodes.get(u));
        }
    });

    return nodes;
}

// Export a graph to a list of edges
function exportGraph(nodes, directed = false) {
    const edges = [];
    const visited = new Set();

    nodes.forEach((node) => {
        node.neighbors.forEach((neighbor) => {
            if (!visited.has(`${neighbor.val}-${node.val}`)) {
                edges.push([node.val, neighbor.val]);
                if (!directed) {
                    visited.add(`${node.val}-${neighbor.val}`);
                }
            }
        });
    });

    return edges;
}

// Exports
module.exports = {
    GraphNode,
    TreeNode,
    ListNode,
    generateTree,
    exportTree,
    generateLinkedList,
    exportLinkedList,
    generateGraph,
    exportGraph,
};