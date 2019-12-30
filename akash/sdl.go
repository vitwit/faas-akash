package akash

import akashTypes "github.com/vitwit/faas-akash/types"

func CreateDefaultDeploymentManifest(image string) akashTypes.SDL {
	return akashTypes.SDL{
		Version: "1.0",
		Services: akashTypes.SDLService{
			Web: akashTypes.SDLWebService{
				Image: image,
				Expose: []akashTypes.SDLWebExpose{
					{
						Port: 80,
						To: []map[string]interface{}{
							{
								"global": true,
							},
						},
					},
				},
			},
		},
		Profiles: akashTypes.SDLProfile{
			Compute: akashTypes.SDLProfileCompute{
				Web: akashTypes.SDLWebCompute{
					CPU:    "0.25",
					Memory: "512Mi",
					Disk:   "1G",
				},
			},
			Placement: akashTypes.SDLProfilePlacement{
				Global: akashTypes.SDLGlobalPlacement{
					Pricing: akashTypes.SDLPricing{
						Web: "1000u",
					},
				},
			},
		},
		Deployment: akashTypes.SDLDeployment{
			Web: akashTypes.SDLDeploymentWeb{
				Global: akashTypes.SDLDeploymentGlobal{
					Profile: "web",
					Count:   1,
				},
			},
		},
	}
}
