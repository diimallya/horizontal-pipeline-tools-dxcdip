package main

const EXPECTEDERROR = "exit status 1"

func PushDbAdminApp(provisionInfo ProvisionInfo) error {
	for userName := range provisionInfo.users {
		err := provisionInfo.provisioner.CallBluemixWithArgs("target -s", userName)
		if err != nil {
			return err
		}

		// Check whether app exists already - if not this call will error
		if err = provisionInfo.provisioner.CallBluemixWithArgs("cf app", removeSpaces(userName)+"-admin"); err != nil {
			if err.Error() == EXPECTEDERROR {
				// Push the app
				if err = provisionInfo.provisioner.CallBluemixWithArgs("cf push -f", "manifest-"+removeSpaces(userName)+".yaml"); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}
