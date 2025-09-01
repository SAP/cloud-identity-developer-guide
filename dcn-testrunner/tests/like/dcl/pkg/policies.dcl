POLICY like_with_underscore {
    GRANT * ON like_with_underscore WHERE s LIKE 'x_y';
}

POLICY like_with_percent {
    GRANT * ON like_with_percent WHERE s LIKE 'x%y';
}

POLICY like_with_escape {
    GRANT * ON like_with_escape WHERE s LIKE 'xö%_ö_y' ESCAPE 'ö';
}
