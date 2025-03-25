package ciscoaciaim

import (
	"fmt"
	ciscoaciaimv1 "github.com/noironetworks/aci-integration-module-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// getVolumes - service volumes
func getVolumes(instance *ciscoaciaimv1v1.CiscoAciAim) []corev1.Volume {
	name := instance.Name
	var scriptsVolumeDefaultMode int32 = 0755
	var config0640AccessMode int32 = 0640

	fernetKeys := []corev1.KeyToPath{}
	numberKeys := int(*instance.Spec.FernetMaxActiveKeys)

	for i := 0; i < numberKeys; i++ {
		fernetKeys = append(
			fernetKeys,
			corev1.KeyToPath{
				Key:  fmt.Sprintf("FernetKeys%d", i),
				Path: fmt.Sprintf("%d", i),
			},
		)
	}

	return []corev1.Volume{
		{
			Name: "scripts",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: &scriptsVolumeDefaultMode,
					SecretName:  name + "-scripts",
				},
			},
		},
		{
			Name: "config-data",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: &config0640AccessMode,
					SecretName:  name + "-config-data",
				},
			},
		},
		{
			Name: "fernet-keys",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: ServiceName,
					Items:      fernetKeys,
				},
			},
		},
		{
			Name: "credential-keys",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: ServiceName,
					Items: []corev1.KeyToPath{
						{
							Key:  "CredentialKeys0",
							Path: "0",
						},
						{
							Key:  "CredentialKeys1",
							Path: "1",
						},
					},
				},
			},
		},
	}

}

func getVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "scripts",
			MountPath: "/usr/local/bin/container-scripts",
			ReadOnly:  true,
		},
		{
			Name:      "config-data",
			MountPath: "/var/lib/config-data/default",
			ReadOnly:  false,
		},
		{
			Name:      "config-data",
			MountPath: "/var/lib/kolla/config_files/config.json",
			SubPath:   "keystone-api-config.json",
			ReadOnly:  true,
		},
		{
			MountPath: "/var/lib/fernet-keys",
			ReadOnly:  true,
			Name:      "fernet-keys",
		},
		{
			MountPath: "/var/lib/credential-keys",
			ReadOnly:  true,
			Name:      "credential-keys",
		},
	}
}
