package indexdb

import "github.com/cdvelop/model"

func (d *indexDB) BackupDataBase(callback func(error)) {
	// reset valores a 0
	d.backups = []backup{}
	d.backupRespond = callback
	d.remaining = 0

	for _, o := range d.objects {
		if len(o.Fields) != 0 { // obtenemos los objetos creados en la db
			d.remaining++
		}
	}

	d.Log("RESPALDANDO BASE DE DATOS INDEX DB")

	// callback(model.Error("NO RESPALDADO AUN"))

	d.addNewObjectsCreated()

}

func (d *indexDB) addNewObjectsCreated() {

	for i, o := range d.objects {
		if len(o.Fields) != 0 {
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
						object:   o,
						data:     data,
						finished: false,
						err:      nil,
					}

					d.backups = append(d.backups, new)

					d.Log("BACKUP REQUERIDO", table)
				}
				d.remaining--

				// finish
				d.finishReadData(index, table)
			})
		}
	}

}

func (d *indexDB) finishReadData(index int, table string) {
	d.Log("INDICE ACTUAL:", index, table)

	d.Log("LECTURA RESTANTE:", d.remaining)

	if d.remaining == 0 {
		d.Log("LECTURA FINALIZADA")
		if len(d.backups) != 0 {
			d.Log("BACKUP A REALZAR:", len(d.backups))
			d.prepareToSendData()

		} else {
			d.Log("BACKUP OK NADA PARA ENVIAR")
			d.backupRespond(nil)
		}
	}
}
