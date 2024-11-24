import ast
import ds_utils

def listy_to_type(stringy_listry_rep, abstract_type):
    """
    Convert a stringy listy representation into the appropriate data structure
    based on the given abstract type.

    Args:
        stringy_listry_rep (str): The stringy representation of the data.
        abstract_type (dict): The abstract type definition with "type" and optional "type_children".

    Returns:
        Any: The converted data structure.
    """
    try:
        # Safely parse the stringy representation into Python structures
        listy_rep = ast.literal_eval(stringy_listry_rep)
    except (ValueError, SyntaxError) as e:
        raise ValueError(f"Failed to parse input: {stringy_listry_rep}. Error: {str(e)}")

    base_type = abstract_type["type"]
    type_children = abstract_type.get("type_children")

    # Handle atomic types
    if base_type in ["Integer", "Double", "String", "Boolean"]:
        # Directly return the value for atomic types
        return listy_rep

    # Handle composite types
    if base_type == "Array":
        # Recursively convert each element using the child type
        return [listy_to_type(repr(item), type_children) for item in listy_rep]

    if base_type == "Matrix":
        # Recursively convert each row using the child type
        return [listy_to_type(repr(row), type_children) for row in listy_rep]

    if base_type == "TreeNode":
        # Use the utility function to generate a TreeNode
        tree = ds_utils.generate_tree(listy_rep)
        if type_children:
            # Convert the value and children using the child type
            tree.val = listy_to_type(repr(tree.val), type_children)
            if tree.left:
                tree.left = listy_to_type(repr(ds_utils.export_tree(tree.left)), abstract_type)
            if tree.right:
                tree.right = listy_to_type(repr(ds_utils.export_tree(tree.right)), abstract_type)
        return tree

    if base_type == "ListNode":
        # Use the utility function to generate a ListNode
        linked_list = ds_utils.generate_linked_list(listy_rep)
        if type_children:
            # Convert each node value using the child type
            current = linked_list
            while current:
                current.val = listy_to_type(repr(current.val), type_children)
                current = current.next
        return linked_list

    if base_type == "Graph":
        # Use the utility function to generate a Graph
        graph = ds_utils.generate_graph(listy_rep)
        if type_children:
            # Convert each node and its neighbors using the child type
            new_adj_list = {}
            for u, neighbors in graph.adj_list.items():
                converted_u = listy_to_type(repr(u), type_children)
                new_adj_list[converted_u] = [listy_to_type(repr(neighbor), type_children) for neighbor in neighbors]
            graph.adj_list = new_adj_list
        return graph

    raise ValueError(f"Unsupported type: {base_type}")
