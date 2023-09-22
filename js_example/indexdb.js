
let db;

function startDataBase() {
    
    db = indexedDB.open("localdb");


    db.addEventListener("error",showDbError);
    db.addEventListener("success",openExistingDB);
    db.addEventListener("upgradeneeded",createNewDB);

    console.log("iniciando base de datos")
}


// startDataBase();


function showDbError(e) {
    
    console.log("showDbError ",e)
}

function openExistingDB(e) {
    
    console.log("openExistingDB ",e)
    
}

function createNewDB(e) {

    console.log("createNewDB ",e)

    let res = e.target.result;

    // crear tabla
    let new_table = res.createObjectStore("table_staff",{keyPath:"id_staff"});

    // crear indices para búsqueda
    new_table.createIndex("searchBYstaff_name","staff_name",{unique:false});
    
}