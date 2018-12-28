const data = [
    ["Contest1", "A+B Problem", "A*B Problem"],
    ["Contest2", "Read File", "Wrie To File"]
];
  
const container = document.getElementById('CompetitionsTable');

const hot = new Handsontable(container, {
    data: data,
    //rowHeaders: true,
    //colHeaders: true,
    filters: true,
    //dropdownMenu: true
});