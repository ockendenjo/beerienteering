# Admin Design

* An index.html page with a linked JavaScript file, and linked CSS file
* The format of the stash file is defined in [stash2025.schema.json](../schema/stash2025.schema.json)

## Functionality

* The page should load JSON from GET `${window.location.host}/stashes` using a fetch request
  * If an error is returned, use an empty StashFile object instead
* The page should save JSON using a PUT request to `${window.location.host}/stashes`
* The page should allow a stash to be added, edited or deleted
  * Don't show the `contents` or `hide` fields in the edit form or the table
  * Label `location` field as `Location description`
  * Don't show the demo mode checkbox
  * Add a note above w3w in the form that `https://www.geocachingtoolbox.com/index.php?lang=en&page=w3w` can be used for conversion
