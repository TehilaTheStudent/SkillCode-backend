from typing import Dict, TypeVar, Generic, List, Optional, Union

T = TypeVar("T")  # Generic type for the value of the node

class TreeNode(Generic[T]):
    def __init__(self, val: T = None, left: Optional['TreeNode[T]'] = None, right: Optional['TreeNode[T]'] = None):
        self.val = val
        self.left = left
        self.right = right

    def __eq__(self, other):
        if not isinstance(other, TreeNode):
            return False
        return (
            self.val == other.val
            and self.left == other.left
            and self.right == other.right
        )

    def __repr__(self):
        return f"TreeNode(val={self.val}, left={self.left}, right={self.right})"

    def __hash__(self):
        return hash((self.val, self.left, self.right))

class ListNode(Generic[T]):
    def __init__(self, val: T = None, next: Optional['ListNode[T]'] = None):
        self.val = val
        self.next = next

    def __eq__(self, other):
        if not isinstance(other, ListNode):
            return False
        current_self, current_other = self, other
        while current_self and current_other:
            if current_self.val != current_other.val:
                return False
            current_self = current_self.next
            current_other = current_other.next
        return current_self is None and current_other is None

    def __repr__(self):
        return f"ListNode(val={self.val}, next={self.next})"

    def __hash__(self):
        current = self
        hash_value = 0
        while current:
            hash_value = hash((hash_value, current.val))
            current = current.next
        return hash_value

class Graph(Generic[T]):
    def __init__(self):
        self.adj_list: Dict[T, List[Union[T, 'Graph[T]', 'TreeNode[T]', 'ListNode[T]']]] = {}

    def add_edge(self, u: T, v: Union[T, 'Graph[T]', 'TreeNode[T]', 'ListNode[T]']):
        if u not in self.adj_list:
            self.adj_list[u] = []
        self.adj_list[u].append(v)

    def __eq__(self, other):
        if not isinstance(other, Graph):
            return False
        if set(self.adj_list.keys()) != set(other.adj_list.keys()):
            return False
        for key in self.adj_list:
            if sorted(self.adj_list[key], key=str) != sorted(other.adj_list[key], key=str):
                return False
        return True

    def __repr__(self):
        return f"Graph(adj_list={self.adj_list})"

    def __hash__(self):
        return hash(frozenset((u, frozenset(neighbors)) for u, neighbors in self.adj_list.items()))

# Generate a tree that supports nested data structures
def generate_tree(values: List[Optional[Union[T, dict]]]) -> Optional[TreeNode[T]]:
    """
    Generate a binary tree from a list of values or nested structures.

    :param values: List of values or dictionaries (e.g., [1, {"val": 2, "left": None}, 3]).
    :return: Root of the constructed TreeNode.
    """
    if not values:
        return None

    root = TreeNode(values[0] if not isinstance(values[0], dict) else TreeNode(**values[0]))
    queue = [root]
    i = 1

    while queue and i < len(values):
        current = queue.pop(0)

        if i < len(values) and values[i] is not None:
            current.left = TreeNode(values[i] if not isinstance(values[i], dict) else TreeNode(**values[i]))
            queue.append(current.left)
        i += 1

        if i < len(values) and values[i] is not None:
            current.right = TreeNode(values[i] if not isinstance(values[i], dict) else TreeNode(**values[i]))
            queue.append(current.right)
        i += 1

    return root

# Generate a graph that supports nested data structures
def generate_graph(edges: List[List[Union[T, dict]]]) -> Graph[T]:
    """
    Generate a graph from a list of directed edges with nested structures.

    :param edges: List of edges (e.g., [[1, 2], [3, {"val": 4}]]).
    :return: A Graph object represented as an adjacency list.
    """
    graph = Graph[T]()
    for edge in edges:
        if len(edge) != 2:
            raise ValueError(f"Invalid edge format: {edge}. Each edge must have exactly two elements.")
        u = edge[0] if not isinstance(edge[0], dict) else TreeNode(**edge[0])
        v = edge[1] if not isinstance(edge[1], dict) else TreeNode(**edge[1])
        graph.add_edge(u, v)
    return graph

# Export tree structure back to a list
def export_tree(root: Optional[TreeNode[T]]) -> List[Optional[Union[T, dict]]]:
    """
    Export a binary tree to a list of values or nested structures.

    :param root: Root of the TreeNode.
    :return: List of values or nested structures.
    """
    if not root:
        return []

    result = []
    queue = [root]

    while queue:
        current = queue.pop(0)
        if current:
            result.append(current.val if not isinstance(current.val, TreeNode) else export_tree(current.val))
            queue.append(current.left)
            queue.append(current.right)
        else:
            result.append(None)

    # Remove trailing None values for cleaner output
    while result and result[-1] is None:
        result.pop()

    return result

# Export graph back to edges
def export_graph(graph: Graph[T]) -> List[List[Union[T, dict]]]:
    """
    Export a graph to a list of directed edges with nested structures.

    :param graph: A Graph object represented as an adjacency list.
    :return: List of edges with nested structures.
    """
    edges = []
    for u, neighbors in graph.adj_list.items():
        for v in neighbors:
            edges.append([u, v])
    return edges

# Generate a linked list from a list of values
def generate_linked_list(values: List[Optional[T]]) -> Optional[ListNode[T]]:
    """
    Generate a singly linked list from a list of values.

    :param values: List of values (e.g., [1, 2, 3]).
    :return: Head of the constructed ListNode.
    """
    if not values:
        return None

    head = ListNode(values[0])
    current = head
    for val in values[1:]:
        current.next = ListNode(val)
        current = current.next

    return head

# Export a linked list to a list of values
def export_linked_list(head: Optional[ListNode[T]]) -> List[Optional[T]]:
    """
    Export a singly linked list to a list of values.

    :param head: Head of the ListNode.
    :return: List of values.
    """
    result = []
    current = head
    while current:
        result.append(current.val)
        current = current.next

    return result
