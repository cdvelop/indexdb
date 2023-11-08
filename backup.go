package indexdb

import "github.com/cdvelop/model"

func (d *indexDB) BackupDataBase(callback func(error)) {
	// reset valores a 0
	d.backups = []backup{}
	d.backupRespond = callback
	d.remaining_reading = len(d.objects)

	d.Log("RESPALDANDO BASE DE DATOS INDEX DB")

	// callback(model.Error("NO RESPALDADO AUN"))

	d.addNewObjectsCreated()

}

func (d *indexDB) addNewObjectsCreated() {

	for i, o := range d.objects {
		index := i // Captura el valor de i en esta iteración
		table := o.Table
		d.ReadAnyDataAsyncInDB(model.ReadDBParams{
			FROM_TABLE:      table,
			WHERE:           []string{"backup"},
			SEARCH_ARGUMENT: "false",
		}, func(data []map[string]interface{}, err error) {
			if err != nil {
				d.Log(err)
				return
			}

			if len(data) != 0 {
				d.Log(data)

				new := backup{
					table:    table,
					data:     data,
					finished: false,
					err:      nil,
				}

				d.backups = append(d.backups, new)

				d.Log("BACKUP REQUERIDO", table)
			}
			d.remaining_reading--

			// finish
			d.finishReadData(index, table)
		})

	}

}

func (d *indexDB) finishReadData(index int, table string) {
	d.Log("INDICE ACTUAL:", index, table)

	d.Log("LECTURA RESTANTE:", d.remaining_reading)

	if d.remaining_reading == 0 {
		d.Log("LECTURA FINALIZADA TOTAL BACKUP A REALZAR:", len(d.backups))
		d.prepareToSendData()
	}
}
