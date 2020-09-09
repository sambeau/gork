# SFSFF: SIMPLE, FLEXIBLE, SAFE-ISH FILE FORMAT

A flexible, tamper-proof file format with file-format, file-format-version amd integrity checking

## ARGUMENTS

    filename
    file-format-identifier
    file-format-version-major
    file-format-version-minor
    secret-key
    []metadata
    data

## WRITE FILE

Pseudocode:

    (file-format-version-major,file-format-version-minor) -> version
    []metadata('\' -> "\\", '\n' -> "\n", '|'' -> "\x7c") ~> join "|" ~-> encoded-metadata
    unshift version onto metadata
    file-identifier + \n + encoded-metadata +\n -> header
    data ~-> gob ~-> zlib-data
    header + secret-key+"1" ~-> header-hash
    header \n header-hash \n zlib-data secret-key+"2" ~-> file-hash
    header \n file-hash \n zlib-data -> OUTFILE
    write OUTFILE

## READ FILE

Pseudocode:

    read to \n -> file-identifier
    check file-identifier ==> ERR
    read to \n -> encoded-metadata
    read to \n -> file-hash
    read to EOF -> zlib data

    integrity check (can skip for faster load) {
        header + secret-key+"1!" ~> header-hash2
        header \n header-hash \n zlib-data secret-key+"2!" ~-> file-hash-2
        check file-hash == file-hash-2 ==> ERR
    }

    encoded-metadata ~=> decodedmetadata ~=> split on `|` ~=> []metadata
    metadata[0] -> (file-format-version-major,file-format-version-minor)
    check file-format-version-major ==> ERR
    zlib-data ~=> gob
    gob ~=> data

## RETURN

    []metadata
    data
