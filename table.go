package indexdb

func (d indexDB) TableExist(table_name string) bool {

	// Obtiene la lista de nombres de almacenes de objetos en la base de datos
	objectStoreNames := d.db.Get("objectStoreNames")
	length := objectStoreNames.Length()

	// Itera a través de los nombres de las tablas y verifica si la tabla ya existe
	for i := 0; i < length; i++ {
		name := objectStoreNames.Index(i).String()
		if name == table_name {
			return true
		}
	}

	return false
}
