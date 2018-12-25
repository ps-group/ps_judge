const data = [
    ["", "A+B Problem", "A*B Problem"],
    ["test_student", 0, 0]
];
  
const container = document.getElementById('ResultsTable');

const hot = new Handsontable(container, {
    data: data,
    rowHeaders: true,
    colHeaders: true,
    filters: true,
    dropdownMenu: true
});