POLICY in_number_array {
    GRANT * ON in_number_array WHERE n IN n_a;
}

POLICY in_boolean_array {
    GRANT * ON in_boolean_array WHERE b IN b_a;
}

POLICY in_string_array {
    GRANT * ON in_string_array WHERE s IN s_a;
}
