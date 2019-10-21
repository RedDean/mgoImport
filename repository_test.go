package mgoImport

import (
	"mgoImport/testUtil"
	"testing"
)

func TestRepository(t *testing.T) {

	//want := map[string]interface{}{
	//    "name":"harden",
	//	"age": 23,
	//	"tel": 10086,
	//}

	t.Run("build properties map", func(t *testing.T) {

		got := &Repository{}
		err := got.buildProperties([]string{"age", "name"}, []string{"int", "string"})
		testUtil.AssertNoError(t, err)

		want := &Repository{
			Properties: []Model{
				{
					FieldType: "int",
					FieldName: "age",
				}, {
					FieldName: "name",
					FieldType: "string",
				},
			},
		}

		testUtil.AssertTwoObjEqual(t, got, want)

	})

	t.Run("build mongo model", func(t *testing.T) {

		repo := &Repository{}
		err1 := repo.buildProperties([]string{"name", "number", "payload"}, []string{"string", "int", "json"})
		testUtil.AssertNoError(t, err1)

		input := []string{"harden", "13", "{\"team\":\"H-town\"}"}

		got, err2 := repo.BuildModel(input)
		testUtil.AssertNoError(t, err2)
		// json 字符串里的数值类型 默认会转成float64，所以map 不相等
		want := map[string]interface{}{
			"name":   "harden",
			"number": 13,
			"team":   "H-town",
		}

		testUtil.AssertTwoObjEqual(t, got, want)
	})

	/*t.Run("insert model", func(t *testing.T) {

		repo := &Repository{}
		err1 := repo.buildProperties([]string{"name", "number", "payload"}, []string{"string", "int", "json"})
		AssertNoError(t, err1)

		input := []string{"harden", "13", "{\"team\":\"H-town\",\"age\":29}"}
		data, err2 := repo.BuildModel(input)
		AssertNoError(t, err2)

		if err := repo.insertData(data); err != nil {
			AssertNoError(t,err)
		}

		got := repo.GetData()


	})*/
	t.Run("test resetBinaries", func(t *testing.T) {
		// arrange
		mockData := map[string]interface{}{
			"apk": map[string]interface{}{
				"href":     "gcs/unity-connect-int/secure/20190521/udp/4c679a19-2fea-4571-9ef6-edd934958255_testDevUnity.apk",
				"fileName": "testDevUnity.apk",
				"fileSize": "22593487",
			},
			"extensions": map[string]interface{}{
				"selectedObfuscationDll": []string{},
				"apkClientId":            "_94lG5fYzcFigLBYy-Afrg",
				"obfuscationDll": []string{
					"UnityEngine.ARModule.dll",
					"UnityEngine.CrashReportingModule.dll",
					"UnityEngine.ParticleSystemModule.dll",
					"UnityEngine.Physics2DModule.dll",
					"UnityEngine.UnityConnectModule.dll",
					"UnityEngine.UnityWebRequestModule.dll",
					"UnityEngine.UnityWebRequestTextureModule.dll",
					"UnityEngine.WindModule.dll",
					"Mono.Security.dll",
					"UnityEngine.AudioModule.dll",
					"UnityEngine.InputModule.dll",
					"UnityEngine.SpriteShapeModule.dll",
					"UnityEngine.SubstanceModule.dll",
					"UnityEngine.TextRenderingModule.dll",
					"UnityEngine.TilemapModule.dll",
					"UnityEngine.UnityWebRequestAudioModule.dll",
					"mscorlib.dll",
					"UnityEngine.Analytics.dll",
					"UnityEngine.TLSModule.dll",
					"UnityEngine.UIModule.dll",
					"UnityEngine.UnityWebRequestWWWModule.dll",
					"UnityEngine.UI.dll",
					"UnityEngine.BaselibModule.dll",
					"UnityEngine.GameCenterModule.dll",
					"UnityEngine.GridModule.dll",
					"UnityEngine.SpatialTracking.dll",
					"System.dll",
					"UnityEngine.AccessibilityModule.dll",
					"UnityEngine.CloudWebServicesModule.dll",
					"UnityEngine.TerrainModule.dll",
					"UnityEngine.UnityAnalyticsModule.dll",
					"Assembly-CSharp-firstpass.dll",
					"UnityEngine.DirectorModule.dll",
					"UnityEngine.UNETModule.dll",
					"System.Core.dll",
					"UnityEngine.CoreModule.dll",
					"UnityEngine.SpriteMaskModule.dll",
					"UnityEngine.AssetBundleModule.dll",
					"UnityEngine.SharedInternalsModule.dll",
					"UnityEngine.StyleSheetsModule.dll",
					"UnityEngine.VRModule.dll",
					"UnityEngine.VehiclesModule.dll",
					"UnityEngine.StandardEvents.dll",
					"UnityEngine.TerrainPhysicsModule.dll",
					"UnityEngine.Timeline.dll",
					"UnityEngine.XRModule.dll",
					"UnityEngine.dll",
					"UnityEngine.AIModule.dll",
					"UnityEngine.AnimationModule.dll",
					"UnityEngine.ClothModule.dll",
					"UnityEngine.ImageConversionModule.dll",
					"UnityEngine.TimelineModule.dll",
					"UnityEngine.HotReloadModule.dll",
					"UnityEngine.Networking.dll",
					"UnityEngine.ParticlesLegacyModule.dll",
					"UnityEngine.PerformanceReportingModule.dll",
					"UnityEngine.PhysicsModule.dll",
					"UnityEngine.ScreenCaptureModule.dll",
					"UnityEngine.SpatialTrackingModule.dll",
					"UnityEngine.UnityWebRequestAssetBundleModule.dll",
					"UnityEngine.VideoModule.dll",
					"UnityEngine.WebModule.dll",
					"UnityEngine.UIElementsModule.dll",
					"UDP.dll",
					"UnityEngine.FacebookModule.dll",
					"UnityEngine.IMGUIModule.dll",
					"UnityEngine.JSONSerializeModule.dll",
					"UnityEngine.UmbraModule.dll",
				},
				"APTOIDE_apks": map[string]interface{}{
					"0.2.10": map[string]interface{}{
						"fileName": "53b33e97-ba45-4b8d-9d53-4951f729a325-signed.apk",
						"fileSize": "22635718",
						"href":     "gcs/unity-connect-int/secure/aptoide/53b33e97-ba45-4b8d-9d53-4951f729a325-signed.apk",
						"type":     "apk",
					},
				},
			},
		}

		// act
		bin := resetBinaries(mockData)

		// assert
		t.Logf("binaries after modified : %v", bin)
		if _, ok := bin.(map[string]interface{})["extensions"]; ok {
			t.Fatalf("still have extensions!")
		}
	})

	t.Run("test resetChannels", func(t *testing.T) {
		// Arrange
		mockData := map[string]interface{}{
			"APTOIDE": map[string]interface{}{
				"extensions": map[string]interface{}{
					"externalItemId":      "com.unity3d.testDevUnity",
					"createStoreTime":     "2019-05-21T11:32Z",
					"publicKey":           "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1joJWm3m0xJoZPlB9omPlj5uMlRYDBL/olDKzYhnrxHuLazdCE9ZYVDYfDUHdmzwPRBykfYsrvoxOHTpxKcfQE3u7reQb3uclGPpqcgdbFpIIX4Q2z7uI/jGgPyp/VhkEUzjM6ErGFBXr92S37ZGxOUg5iYTi7r8z4tiGzyJrBjQjp/TWbzhIRpXELea6XgzqDRwyS/fP2gPceE4E2PjXon8LCaioBaiY61Sz5gs9er+VSgh7mP3zXKf9cLt8n6OqzjM4n6FW5NZtVjQOnTxZ/OjIQlxJnfzNsNvVA21tXPUK2HGZZQfNvRxrOvz4KwdgIqUPp1xXcODJkWIsyvpywIDAQAB",
					"packageName":         "com.unity3d.testDevUnity",
					"targetSdkId":         "274877919251",
					"targetStep":          "PRODUCTION",
					"packingStatus":       "COMPLETED",
					"packedPackageName":   "com.unity3d.testDevUnity",
					"packingCompleteTime": "2019-05-21T11:33:06Z",
					"packingId":           "12369505849644",
					"syncStatus":          "ERROR",
					"syncTime":            "2019-05-21T11:37:40Z",
				},
			},
		}

		want := map[string]interface{}{
			"APTOIDE": map[string]interface{}{
				"externalItemId":      "com.unity3d.testDevUnity",
				"createStoreTime":     "2019-05-21T11:32Z",
				"publicKey":           "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1joJWm3m0xJoZPlB9omPlj5uMlRYDBL/olDKzYhnrxHuLazdCE9ZYVDYfDUHdmzwPRBykfYsrvoxOHTpxKcfQE3u7reQb3uclGPpqcgdbFpIIX4Q2z7uI/jGgPyp/VhkEUzjM6ErGFBXr92S37ZGxOUg5iYTi7r8z4tiGzyJrBjQjp/TWbzhIRpXELea6XgzqDRwyS/fP2gPceE4E2PjXon8LCaioBaiY61Sz5gs9er+VSgh7mP3zXKf9cLt8n6OqzjM4n6FW5NZtVjQOnTxZ/OjIQlxJnfzNsNvVA21tXPUK2HGZZQfNvRxrOvz4KwdgIqUPp1xXcODJkWIsyvpywIDAQAB",
				"packageName":         "com.unity3d.testDevUnity",
				"targetSdkId":         "274877919251",
				"targetStep":          "PRODUCTION",
				"packingStatus":       "COMPLETED",
				"packedPackageName":   "com.unity3d.testDevUnity",
				"packingCompleteTime": "2019-05-21T11:33:06Z",
				"packingId":           "12369505849644",
				"syncStatus":          "ERROR",
				"syncTime":            "2019-05-21T11:37:40Z",
			},
		}

		// Act
		got := resetChannels(mockData)

		// Assert
		testUtil.AssertTwoObjEqual(t, got, want)
	})
}
