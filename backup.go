package indexdb

import (
	"github.com/cdvelop/model"
)

func (d *indexDB) BackupDataBase(callback func(err string)) {
	// reset valores a 0
	d.backups = []backup{}
	d.backupRespond = callback
	d.remaining = len(d.getObjectsDB())

	d.Log("RESPALDANDO BASE DE DATOS INDEX DB")

	// callback("NO RESPALDADO AUN"))

	d.addNewObjectsCreated()

}

func (d *indexDB) addNewObjectsCreated() {

	for i, o := range d.getObjectsDB() {
		index := i // Captura el valor de i en esta iteración
		table := o.Table
		d.ReadAsyncDataDB(model.ReadParams{
			FROM_TABLE: table,
			WHERE:      []map[string]string{{"create": "backup", "update": "backup", "delete": "backup"}},
			RETURN_ANY: true,
		}, func(r *model.ReadResults, err string) {

			if err != "" {
				d.Log(err)
				return
			}

			if len(r.ResultsAny) != 0 {
				d.Log(r.ResultsAny)

				new := backup{
					object:   o,
					data:     r.ResultsAny,
					finished: false,
					err:      "",
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
			d.backupRespond("")
		}
	}
}
