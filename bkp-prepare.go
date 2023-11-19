package indexdb

func (d *indexDB) prepareToSendData() {

	for _, b := range d.backups {

		d.Log("OBJETO", b.object.ObjectName)
		// var o *model.Object
		// for _, v := range d.objects {
		// 	if v.Table == o.Table {
		// 		o = v
		// 		break
		// 	}
		// }

		for _, item := range b.data {

			if b.object.Table == "file" {
				d.Log("TIPO FILE ENVIÓ FORM DATA", item)

				// d.http.SendAllRequests()

			} else {
				d.Log("ENVIÓ NORMAL JSON", item)

				if _, update := item["update"]; update {

					d.Log("UPDATE", item)

				} else {

					d.Log("CREATE", item)

				}
			}

		}

	}

}
