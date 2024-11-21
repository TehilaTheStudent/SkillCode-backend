from typing import TypeVar, Generic, List

T = TypeVar('T')  # Generic type for the value of the node

class Graph(Generic[T]):
    def __init__(self, val: T = None, neighbors: List['Graph[T]'] = None):
        """
        A node in a graph.

        :param val: The value of the node (any type specified by T).
        :param neighbors: List of neighbor nodes (Graph).
        """
        self.val = val
        self.neighbors = neighbors if neighbors is not None else []

    def __repr__(self):
        return f"Graph(val={self.val})"

from typing import Optional

class TreeNode(Generic[T]):
    def __init__(self, val: T = None, left: Optional['TreeNode[T]'] = None, right: Optional['TreeNode[T]'] = None):
        """
        A node in a binary tree.

        :param val: The value of the node (any type specified by T).
        :param left: Left child (TreeNode).
        :param right: Right child (TreeNode).
        """
        self.val = val
        self.left = left
        self.right = right

    def __repr__(self):
        return f"TreeNode(val={self.val})"

class ListNode(Generic[T]):
    def __init__(self, val: T = None, next: Optional['ListNode[T]'] = None):
        """
        A node in a singly linked list.

        :param val: The value of the node (any type specified by T).
        :param next: The next node (ListNode).
        """
        self.val = val
        self.next = next

    def __repr__(self):
        return f"ListNode(val={self.val})"

def generate_tree(values: List[Optional[T]]) -> Optional[TreeNode[T]]:
    """
    Generate a binary tree from a list of values (BFS style).

    :param values: List of values (e.g., [1, 2, 3, None, None, 4, 5]).
    :return: Root of the constructed TreeNode.
    """
    if not values:
        return None

    root = TreeNode(values[0])
    queue = [root]
    i = 1

    while queue and i < len(values):
        current = queue.pop(0)

        if values[i] is not None:
            current.left = TreeNode(values[i])
            queue.append(current.left)
        i += 1

        if i < len(values) and values[i] is not None:
            current.right = TreeNode(values[i])
            queue.append(current.right)
        i += 1

    return root

def export_tree(root: Optional[TreeNode[T]]) -> List[Optional[T]]:
    """
    Export a binary tree to a list of values (BFS style).

    :param root: Root of the TreeNode.
    :return: List of values (e.g., [1, 2, 3, None, None, 4, 5]).
    """
    if not root:
        return []

    result = []
    queue = [root]

    while queue:
        current = queue.pop(0)
        if current:
            result.append(current.val)
            queue.append(current.left)
            queue.append(current.right)
        else:
            result.append(None)

    # Remove trailing None values for cleaner output
    while result and result[-1] is None:
        result.pop()

    return result

def generate_linked_list(values: List[T]) -> Optional[ListNode[T]]:
    """
    Generate a singly linked list from a list of values.

    :param values: List of values (e.g., [1, 2, 3, 4]).
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

def export_linked_list(head: Optional[ListNode[T]]) -> List[T]:
    """
    Export a singly linked list to a list of values.

    :param head: Head of the ListNode.
    :return: List of values (e.g., [1, 2, 3, 4]).
    """
    result = []
    current = head

    while current:
        result.append(current.val)
        current = current.next

    return result

def generate_graph(edges: List[List[T]], directed: bool = False) -> dict[T, Graph[T]]:
    """
    Generate a graph from a list of edges.

    :param edges: List of edges (e.g., [[1, 2], [1, 3], [4, 2]]).
    :param directed: Boolean indicating if the graph is directed.
    :return: Dictionary of Graph objects representing the graph.
    """
    nodes = {}

    for edge in edges:
        u, v = edge

        if u not in nodes:
            nodes[u] = Graph(u)
        if v not in nodes:
            nodes[v] = Graph(v)

        nodes[u].neighbors.append(nodes[v])
        if not directed:
            nodes[v].neighbors.append(nodes[u])

    return nodes

def export_graph(nodes: dict[T, Graph[T]], directed: bool = False) -> List[List[T]]:
    """
    Export a graph to a list of edges.

    :param nodes: Dictionary of Graph objects representing the graph.
    :return: List of edges (e.g., [[1, 2], [1, 3], [4, 2]]).
    """
    edges = []
    visited = set()

    for node in nodes.values():
        for neighbor in node.neighbors:
            if (node.val, neighbor.val) not in visited:
                edges.append([node.val, neighbor.val])
                if not directed:  # Prevent duplicate undirected edges
                    visited.add((neighbor.val, node.val))

    return edges
