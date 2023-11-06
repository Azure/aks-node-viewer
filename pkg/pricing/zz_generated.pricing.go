//go:build !ignore_autogenerated

// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package pricing

import "time"

// generated at 2023-03-01T19:26:21Z for eastus

var initialPriceUpdate, _ = time.Parse(time.RFC3339, "2023-03-01T19:26:21Z")
var initialOnDemandPrices = map[string]map[string]float64{}

func init() {
	// eastus
	initialOnDemandPrices["eastus"] = map[string]float64{
		"Basic_A0":                  0.018000,
		"Basic_A1":                  0.023000,
		"Basic_A2":                  0.079000,
		"Basic_A3":                  0.176000,
		"Basic_A4":                  0.352000,
		"DCadsv5 Type 1":            6.979000,
		"DCasv5 Type 1":             5.827000,
		"DCdsv3 Type1":              5.966000,
		"DCsv2 Type 1":              0.845000,
		"DCsv3 Type1":               5.069000,
		"Dadsv5_Type1":              6.345000,
		"Dasv4_Type1":               5.069000,
		"Dasv4_Type2":               5.280000,
		"Dasv5_Type1":               5.298000,
		"Ddsv4_Type 1":              3.978000,
		"Ddsv4_Type2":               4.723000,
		"Ddsv5_Type1":               5.966000,
		"Dsv3_Type1":                3.380000,
		"Dsv3_Type2":                3.802000,
		"Dsv3_Type3":                4.225000,
		"Dsv3_Type4":                5.281000,
		"Dsv4_Type1":                4.224000,
		"Dsv4_Type2":                5.069000,
		"Dsv5_Type1":                5.280000,
		"ECadsv5 Type 1":            8.369000,
		"ECasv5 Type 1":             7.219000,
		"Eadsv5_Type1":              6.917000,
		"Easv4_Type1":               6.653000,
		"Easv4_Type2":               6.653000,
		"Easv5_Type1":               5.966000,
		"Ebdsv5-Type1":              5.878000,
		"Ebsv5-Type1":               5.245000,
		"Edsv4_Type 1":              5.069000,
		"Edsv4_Type2":               6.019000,
		"Edsv5_Type1":               7.603000,
		"Esv3_Type1":                3.881000,
		"Esv3_Type2":                4.297000,
		"Esv3_Type3":                5.545000,
		"Esv3_Type4":                5.822000,
		"Esv4_Type1":                4.435000,
		"Esv4_Type2":                5.821000,
		"Esv5_Type1":                6.653000,
		"FXmds Type1":               4.910000,
		"Fsv2 Type3":                3.927000,
		"Fsv2_Type2":                3.366000,
		"Fsv2_Type4":                4.468000,
		"Lasv3_Type1":               6.864000,
		"Lsv2_Type1":                6.864000,
		"Lsv3_Type1":                7.656000,
		"Mdsmv2MedMem _Type1":       29.359000,
		"Mdsv2MedMem_Type1":         14.674000,
		"Ms_Type1":                  14.669000,
		"Msm_Type1":                 29.363000,
		"Msmv2MedMem Type1":         28.896000,
		"Msmv2_Type1":               109.060000,
		"Msv2MedMem Type1":          14.211000,
		"Msv2_Type1":                54.538000,
		"NVasv4_Type1":              7.176000,
		"NVsv3_Type1":               5.016000,
		"Standard_A0":               0.020000,
		"Standard_A1":               0.060000,
		"Standard_A1_v2":            0.043000,
		"Standard_A2":               0.120000,
		"Standard_A2_v2":            0.091000,
		"Standard_A2m_v2":           0.119000,
		"Standard_A3":               0.240000,
		"Standard_A4":               0.480000,
		"Standard_A4_v2":            0.191000,
		"Standard_A4m_v2":           0.238000,
		"Standard_A5":               0.250000,
		"Standard_A6":               0.500000,
		"Standard_A7":               1.000000,
		"Standard_A8_v2":            0.400000,
		"Standard_A8m_v2":           0.475000,
		"Standard_B12ms":            0.499000,
		"Standard_B16als_v2":        0.567000,
		"Standard_B16as_v2":         0.640000,
		"Standard_B16ls_v2":         0.589000,
		"Standard_B16ms":            0.666000,
		"Standard_B16pls_v2":        0.425000,
		"Standard_B16ps_v2":         0.480000,
		"Standard_B16s_v2":          0.666000,
		"Standard_B1ls":             0.005200,
		"Standard_B1ms":             0.020700,
		"Standard_B1s":              0.010400,
		"Standard_B20ms":            0.832000,
		"Standard_B2als_v2":         0.040000,
		"Standard_B2as_v2":          0.080000,
		"Standard_B2ats_v2":         0.010000,
		"Standard_B2ls_v2":          0.041600,
		"Standard_B2ms":             0.083200,
		"Standard_B2pls_v2":         0.030000,
		"Standard_B2ps_v2":          0.060000,
		"Standard_B2pts_v2":         0.007500,
		"Standard_B2s":              0.041600,
		"Standard_B2s_v2":           0.083200,
		"Standard_B2ts_v2":          0.010400,
		"Standard_B32als_v2":        1.133000,
		"Standard_B32as_v2":         1.280000,
		"Standard_B32ls_v2":         1.179000,
		"Standard_B32s_v2":          1.331000,
		"Standard_B4als_v2":         0.142000,
		"Standard_B4as_v2":          0.160000,
		"Standard_B4ls_v2":          0.147000,
		"Standard_B4ms":             0.166000,
		"Standard_B4pls_v2":         0.106000,
		"Standard_B4ps_v2":          0.120000,
		"Standard_B4s_v2":           0.016600,
		"Standard_B8als_v2":         0.283000,
		"Standard_B8as_v2":          0.320000,
		"Standard_B8ls_v2":          0.295000,
		"Standard_B8ms":             0.333000,
		"Standard_B8pls_v2":         0.213000,
		"Standard_B8ps_v2":          0.240000,
		"Standard_B8s_v2":           0.333000,
		"Standard_D1":               0.077000,
		"Standard_D11":              0.193000,
		"Standard_D11_v2":           0.185000,
		"Standard_D11_v2_Promo":     0.185000,
		"Standard_D12":              0.386000,
		"Standard_D12_v2":           0.371000,
		"Standard_D12_v2_Promo":     0.371000,
		"Standard_D13":              0.771000,
		"Standard_D13_v2":           0.741000,
		"Standard_D13_v2_Promo":     0.741000,
		"Standard_D14":              1.542000,
		"Standard_D14_v2":           1.482000,
		"Standard_D14_v2_Promo":     1.482000,
		"Standard_D15_v2":           1.853000,
		"Standard_D15i_v2":          1.853000,
		"Standard_D16_v3":           0.768000,
		"Standard_D16_v4":           0.768000,
		"Standard_D16_v5":           0.768000,
		"Standard_D16a_v4":          0.768000,
		"Standard_D16ads_v5":        0.824000,
		"Standard_D16as_v4":         0.768000,
		"Standard_D16as_v5":         0.688000,
		"Standard_D16d_v4":          0.904000,
		"Standard_D16d_v5":          0.904000,
		"Standard_D16ds_v4":         0.904000,
		"Standard_D16ds_v5":         0.904000,
		"Standard_D16lds_v5":        0.768000,
		"Standard_D16ls_v5":         0.680000,
		"Standard_D16pds_v5":        0.723000,
		"Standard_D16plds_v5":       0.614000,
		"Standard_D16pls_v5":        0.544000,
		"Standard_D16ps_v5":         0.616000,
		"Standard_D16s_v3":          0.768000,
		"Standard_D16s_v4":          0.768000,
		"Standard_D16s_v5":          0.768000,
		"Standard_D1_v2":            0.073000,
		"Standard_D2":               0.154000,
		"Standard_D2_v2":            0.146000,
		"Standard_D2_v2_Promo":      0.146000,
		"Standard_D2_v3":            0.096000,
		"Standard_D2_v4":            0.096000,
		"Standard_D2_v5":            0.096000,
		"Standard_D2a_v4":           0.096000,
		"Standard_D2ads_v5":         0.103000,
		"Standard_D2as_v4":          0.096000,
		"Standard_D2as_v5":          0.086000,
		"Standard_D2d_v4":           0.113000,
		"Standard_D2d_v5":           0.113000,
		"Standard_D2ds_v4":          0.113000,
		"Standard_D2ds_v5":          0.113000,
		"Standard_D2lds_v5":         0.096000,
		"Standard_D2ls_v5":          0.085000,
		"Standard_D2pds_v5":         0.090400,
		"Standard_D2plds_v5":        0.076800,
		"Standard_D2pls_v5":         0.068000,
		"Standard_D2ps_v5":          0.077000,
		"Standard_D2s_v3":           0.096000,
		"Standard_D2s_v4":           0.096000,
		"Standard_D2s_v5":           0.096000,
		"Standard_D3":               0.308000,
		"Standard_D32-16s_v3":       1.536000,
		"Standard_D32-8s_v3":        1.536000,
		"Standard_D32_v3":           1.536000,
		"Standard_D32_v4":           1.536000,
		"Standard_D32_v5":           1.536000,
		"Standard_D32a_v4":          1.536000,
		"Standard_D32ads_v5":        1.648000,
		"Standard_D32as_v4":         1.536000,
		"Standard_D32as_v5":         1.376000,
		"Standard_D32d_v4":          1.808000,
		"Standard_D32d_v5":          1.808000,
		"Standard_D32ds_v4":         1.808000,
		"Standard_D32ds_v5":         1.808000,
		"Standard_D32lds_v5":        1.536000,
		"Standard_D32ls_v5":         1.360000,
		"Standard_D32pds_v5":        1.446000,
		"Standard_D32plds_v5":       1.229000,
		"Standard_D32pls_v5":        1.088000,
		"Standard_D32ps_v5":         1.232000,
		"Standard_D32s_v3":          1.536000,
		"Standard_D32s_v4":          1.536000,
		"Standard_D32s_v5":          1.536000,
		"Standard_D3_v2":            0.293000,
		"Standard_D3_v2_Promo":      0.293000,
		"Standard_D4":               0.616000,
		"Standard_D48_v3":           2.304000,
		"Standard_D48_v4":           2.304000,
		"Standard_D48_v5":           2.304000,
		"Standard_D48a_v4":          2.304000,
		"Standard_D48ads_v5":        2.472000,
		"Standard_D48as_v4":         2.304000,
		"Standard_D48as_v5":         2.064000,
		"Standard_D48d_v4":          2.712000,
		"Standard_D48d_v5":          2.712000,
		"Standard_D48ds_v4":         2.712000,
		"Standard_D48ds_v5":         2.712000,
		"Standard_D48lds_v5":        2.304000,
		"Standard_D48ls_v5":         2.040000,
		"Standard_D48pds_v5":        2.170000,
		"Standard_D48plds_v5":       1.843000,
		"Standard_D48pls_v5":        1.632000,
		"Standard_D48ps_v5":         1.848000,
		"Standard_D48s_v3":          2.304000,
		"Standard_D48s_v4":          2.304000,
		"Standard_D48s_v5":          2.304000,
		"Standard_D4_v2":            0.585000,
		"Standard_D4_v2_Promo":      0.585000,
		"Standard_D4_v3":            0.192000,
		"Standard_D4_v4":            0.192000,
		"Standard_D4_v5":            0.192000,
		"Standard_D4a_v4":           0.192000,
		"Standard_D4ads_v5":         0.206000,
		"Standard_D4as_v4":          0.192000,
		"Standard_D4as_v5":          0.172000,
		"Standard_D4d_v4":           0.226000,
		"Standard_D4d_v5":           0.226000,
		"Standard_D4ds_v4":          0.226000,
		"Standard_D4ds_v5":          0.226000,
		"Standard_D4lds_v5":         0.192000,
		"Standard_D4ls_v5":          0.170000,
		"Standard_D4pds_v5":         0.181000,
		"Standard_D4plds_v5":        0.154000,
		"Standard_D4pls_v5":         0.136000,
		"Standard_D4ps_v5":          0.154000,
		"Standard_D4s_v3":           0.192000,
		"Standard_D4s_v4":           0.192000,
		"Standard_D4s_v5":           0.192000,
		"Standard_D5_v2":            1.170000,
		"Standard_D5_v2_Promo":      1.170000,
		"Standard_D64-16s_v3":       3.072000,
		"Standard_D64-32s_v3":       3.072000,
		"Standard_D64_v3":           3.072000,
		"Standard_D64_v4":           3.072000,
		"Standard_D64_v5":           3.072000,
		"Standard_D64a_v4":          3.072000,
		"Standard_D64ads_v5":        3.296000,
		"Standard_D64as_v4":         3.072000,
		"Standard_D64as_v5":         2.752000,
		"Standard_D64d_v4":          3.616000,
		"Standard_D64d_v5":          3.616000,
		"Standard_D64ds_v4":         3.616000,
		"Standard_D64ds_v5":         3.616000,
		"Standard_D64lds_v5":        3.072000,
		"Standard_D64ls_v5":         2.720000,
		"Standard_D64pds_v5":        2.893000,
		"Standard_D64plds_v5":       2.458000,
		"Standard_D64pls_v5":        2.176000,
		"Standard_D64ps_v5":         2.464000,
		"Standard_D64s_v3":          3.072000,
		"Standard_D64s_v4":          3.072000,
		"Standard_D64s_v5":          3.072000,
		"Standard_D8_v3":            0.384000,
		"Standard_D8_v4":            0.384000,
		"Standard_D8_v5":            0.384000,
		"Standard_D8a_v4":           0.384000,
		"Standard_D8ads_v5":         0.412000,
		"Standard_D8as_v4":          0.384000,
		"Standard_D8as_v5":          0.344000,
		"Standard_D8d_v4":           0.452000,
		"Standard_D8d_v5":           0.452000,
		"Standard_D8ds_v4":          0.452000,
		"Standard_D8ds_v5":          0.452000,
		"Standard_D8lds_v5":         0.384000,
		"Standard_D8ls_v5":          0.340000,
		"Standard_D8pds_v5":         0.362000,
		"Standard_D8plds_v5":        0.307000,
		"Standard_D8pls_v5":         0.272000,
		"Standard_D8ps_v5":          0.308000,
		"Standard_D8s_v3":           0.384000,
		"Standard_D8s_v4":           0.384000,
		"Standard_D8s_v5":           0.384000,
		"Standard_D96_v5":           4.608000,
		"Standard_D96a_v4":          4.608000,
		"Standard_D96ads_v5":        4.944000,
		"Standard_D96as_v4":         4.608000,
		"Standard_D96as_v5":         4.128000,
		"Standard_D96d_v5":          5.424000,
		"Standard_D96ds_v5":         5.424000,
		"Standard_D96lds_v5":        4.608000,
		"Standard_D96ls_v5":         4.080000,
		"Standard_D96s_v5":          4.608000,
		"Standard_DC16ads_v5":       0.906000,
		"Standard_DC16as_v5":        0.757000,
		"Standard_DC16ds_v3":        1.808000,
		"Standard_DC16s_v3":         1.536000,
		"Standard_DC1ds_v3":         0.113000,
		"Standard_DC1s_v2":          0.096000,
		"Standard_DC1s_v3":          0.096000,
		"Standard_DC24ds_v3":        2.712000,
		"Standard_DC24s_v3":         2.304000,
		"Standard_DC2ads_v5":        0.113000,
		"Standard_DC2as_v5":         0.094600,
		"Standard_DC2ds_v3":         0.226000,
		"Standard_DC2s":             0.198000,
		"Standard_DC2s_v2":          0.192000,
		"Standard_DC2s_v3":          0.192000,
		"Standard_DC32ads_v5":       1.813000,
		"Standard_DC32as_v5":        1.514000,
		"Standard_DC32ds_v3":        3.616000,
		"Standard_DC32s_v3":         3.072000,
		"Standard_DC48ads_v5":       2.719000,
		"Standard_DC48as_v5":        2.270000,
		"Standard_DC48ds_v3":        5.424000,
		"Standard_DC48s_v3":         4.608000,
		"Standard_DC4ads_v5":        0.227000,
		"Standard_DC4as_v5":         0.189000,
		"Standard_DC4ds_v3":         0.452000,
		"Standard_DC4s":             0.395000,
		"Standard_DC4s_v2":          0.384000,
		"Standard_DC4s_v3":          0.384000,
		"Standard_DC64ads_v5":       3.626000,
		"Standard_DC64as_v5":        3.027000,
		"Standard_DC8_v2":           0.768000,
		"Standard_DC8ads_v5":        0.453000,
		"Standard_DC8as_v5":         0.378000,
		"Standard_DC8ds_v3":         0.904000,
		"Standard_DC8s_v3":          0.768000,
		"Standard_DC96ads_v5":       5.438000,
		"Standard_DC96as_v5":        4.541000,
		"Standard_DS1":              0.077000,
		"Standard_DS11":             0.193000,
		"Standard_DS11-1_v2":        0.185000,
		"Standard_DS11_v2":          0.185000,
		"Standard_DS11_v2_Promo":    0.185000,
		"Standard_DS12":             0.386000,
		"Standard_DS12-1_v2":        0.371000,
		"Standard_DS12-2_v2":        0.371000,
		"Standard_DS12_v2":          0.371000,
		"Standard_DS12_v2_Promo":    0.371000,
		"Standard_DS13":             0.771000,
		"Standard_DS13-2_v2":        0.741000,
		"Standard_DS13-4_v2":        0.741000,
		"Standard_DS13_v2":          0.741000,
		"Standard_DS13_v2_Promo":    0.741000,
		"Standard_DS14":             1.542000,
		"Standard_DS14-4_v2":        1.482000,
		"Standard_DS14-8_v2":        1.482000,
		"Standard_DS14_v2":          1.482000,
		"Standard_DS14_v2_Promo":    1.482000,
		"Standard_DS15_v2":          1.853000,
		"Standard_DS15i_v2":         1.853000,
		"Standard_DS1_v2":           0.073000,
		"Standard_DS2":              0.154000,
		"Standard_DS2_v2":           0.146000,
		"Standard_DS2_v2_Promo":     0.146000,
		"Standard_DS3":              0.308000,
		"Standard_DS3_v2":           0.293000,
		"Standard_DS3_v2_Promo":     0.293000,
		"Standard_DS4":              0.616000,
		"Standard_DS4_v2":           0.585000,
		"Standard_DS4_v2_Promo":     0.585000,
		"Standard_DS5_v2":           1.170000,
		"Standard_DS5_v2_Promo":     1.170000,
		"Standard_E104i_v5":         7.207000,
		"Standard_E104id_v5":        8.237000,
		"Standard_E104ids_v5":       8.237000,
		"Standard_E104is_v5":        7.207000,
		"Standard_E112iads_v5":      8.070000,
		"Standard_E112ias_v5":       6.961000,
		"Standard_E112ibds_v5":      10.287000,
		"Standard_E112ibs_v5":       9.178000,
		"Standard_E16-4ads_v5":      1.048000,
		"Standard_E16-4as_v4":       1.008000,
		"Standard_E16-4as_v5":       0.904000,
		"Standard_E16-4ds_v4":       1.152000,
		"Standard_E16-4ds_v5":       1.152000,
		"Standard_E16-4s_v3":        1.008000,
		"Standard_E16-4s_v4":        1.008000,
		"Standard_E16-4s_v5":        1.008000,
		"Standard_E16-8ads_v5":      1.048000,
		"Standard_E16-8as_v4":       1.008000,
		"Standard_E16-8as_v5":       0.904000,
		"Standard_E16-8ds_v4":       1.152000,
		"Standard_E16-8ds_v5":       1.152000,
		"Standard_E16-8s_v3":        1.008000,
		"Standard_E16-8s_v4":        1.008000,
		"Standard_E16-8s_v5":        1.008000,
		"Standard_E16_v3":           1.008000,
		"Standard_E16_v4":           1.008000,
		"Standard_E16_v5":           1.008000,
		"Standard_E16a_v4":          1.008000,
		"Standard_E16ads_v5":        1.048000,
		"Standard_E16as_v4":         1.008000,
		"Standard_E16as_v5":         0.904000,
		"Standard_E16bds_v5":        1.336000,
		"Standard_E16bs_v5":         1.192000,
		"Standard_E16d_v4":          1.152000,
		"Standard_E16d_v5":          1.152000,
		"Standard_E16ds_v4":         1.152000,
		"Standard_E16ds_v5":         1.152000,
		"Standard_E16pds_v5":        0.922000,
		"Standard_E16ps_v5":         0.806000,
		"Standard_E16s_v3":          1.008000,
		"Standard_E16s_v4":          1.008000,
		"Standard_E16s_v5":          1.008000,
		"Standard_E20_v3":           1.260000,
		"Standard_E20_v4":           1.260000,
		"Standard_E20_v5":           1.260000,
		"Standard_E20a_v4":          1.260000,
		"Standard_E20ads_v5":        1.310000,
		"Standard_E20as_v4":         1.260000,
		"Standard_E20as_v5":         1.130000,
		"Standard_E20d_v4":          1.440000,
		"Standard_E20d_v5":          1.440000,
		"Standard_E20ds_v4":         1.440000,
		"Standard_E20ds_v5":         1.440000,
		"Standard_E20pds_v5":        1.152000,
		"Standard_E20ps_v5":         1.008000,
		"Standard_E20s_v3":          1.260000,
		"Standard_E20s_v4":          1.260000,
		"Standard_E20s_v5":          1.260000,
		"Standard_E2_v3":            0.126000,
		"Standard_E2_v4":            0.126000,
		"Standard_E2_v5":            0.126000,
		"Standard_E2a_v4":           0.126000,
		"Standard_E2ads_v5":         0.131000,
		"Standard_E2as_v4":          0.126000,
		"Standard_E2as_v5":          0.113000,
		"Standard_E2bds_v5":         0.167000,
		"Standard_E2bs_v5":          0.149000,
		"Standard_E2d_v4":           0.144000,
		"Standard_E2d_v5":           0.144000,
		"Standard_E2ds_v4":          0.144000,
		"Standard_E2ds_v5":          0.144000,
		"Standard_E2pds_v5":         0.115000,
		"Standard_E2ps_v5":          0.101000,
		"Standard_E2s_v3":           0.126000,
		"Standard_E2s_v4":           0.126000,
		"Standard_E2s_v5":           0.126000,
		"Standard_E32-16ads_v5":     2.096000,
		"Standard_E32-16as_v4":      2.016000,
		"Standard_E32-16as_v5":      1.808000,
		"Standard_E32-16ds_v4":      2.304000,
		"Standard_E32-16ds_v5":      2.304000,
		"Standard_E32-16s_v3":       2.016000,
		"Standard_E32-16s_v4":       2.016000,
		"Standard_E32-16s_v5":       2.016000,
		"Standard_E32-8ads_v5":      2.096000,
		"Standard_E32-8as_v4":       2.016000,
		"Standard_E32-8as_v5":       1.808000,
		"Standard_E32-8ds_v4":       2.304000,
		"Standard_E32-8ds_v5":       2.304000,
		"Standard_E32-8s_v3":        2.016000,
		"Standard_E32-8s_v4":        2.016000,
		"Standard_E32-8s_v5":        2.016000,
		"Standard_E32_v3":           2.016000,
		"Standard_E32_v4":           2.016000,
		"Standard_E32_v5":           2.016000,
		"Standard_E32a_v4":          2.016000,
		"Standard_E32ads_v5":        2.096000,
		"Standard_E32as_v4":         2.016000,
		"Standard_E32as_v5":         1.808000,
		"Standard_E32bds_v5":        2.672000,
		"Standard_E32bs_v5":         2.384000,
		"Standard_E32d_v4":          2.304000,
		"Standard_E32d_v5":          2.304000,
		"Standard_E32ds_v4":         2.304000,
		"Standard_E32ds_v5":         2.304000,
		"Standard_E32pds_v5":        1.843000,
		"Standard_E32ps_v5":         1.613000,
		"Standard_E32s_v3":          2.016000,
		"Standard_E32s_v4":          2.016000,
		"Standard_E32s_v5":          2.016000,
		"Standard_E4-2ads_v5":       0.262000,
		"Standard_E4-2as_v4":        0.252000,
		"Standard_E4-2as_v5":        0.226000,
		"Standard_E4-2ds_v4":        0.288000,
		"Standard_E4-2ds_v5":        0.288000,
		"Standard_E4-2s_v3":         0.252000,
		"Standard_E4-2s_v4":         0.252000,
		"Standard_E4-2s_v5":         0.252000,
		"Standard_E48_v3":           3.024000,
		"Standard_E48_v4":           3.024000,
		"Standard_E48_v5":           3.024000,
		"Standard_E48a_v4":          3.024000,
		"Standard_E48ads_v5":        3.144000,
		"Standard_E48as_v4":         3.024000,
		"Standard_E48as_v5":         2.712000,
		"Standard_E48bds_v5":        4.008000,
		"Standard_E48bs_v5":         3.576000,
		"Standard_E48d_v4":          3.456000,
		"Standard_E48d_v5":          3.456000,
		"Standard_E48ds_v4":         3.456000,
		"Standard_E48ds_v5":         3.456000,
		"Standard_E48s_v3":          3.024000,
		"Standard_E48s_v4":          3.024000,
		"Standard_E48s_v5":          3.024000,
		"Standard_E4_v3":            0.252000,
		"Standard_E4_v4":            0.252000,
		"Standard_E4_v5":            0.252000,
		"Standard_E4a_v4":           0.252000,
		"Standard_E4ads_v5":         0.262000,
		"Standard_E4as_v4":          0.252000,
		"Standard_E4as_v5":          0.226000,
		"Standard_E4bds_v5":         0.334000,
		"Standard_E4bs_v5":          0.298000,
		"Standard_E4d_v4":           0.288000,
		"Standard_E4d_v5":           0.288000,
		"Standard_E4ds_v4":          0.288000,
		"Standard_E4ds_v5":          0.288000,
		"Standard_E4pds_v5":         0.230000,
		"Standard_E4ps_v5":          0.202000,
		"Standard_E4s_v3":           0.252000,
		"Standard_E4s_v4":           0.252000,
		"Standard_E4s_v5":           0.252000,
		"Standard_E64-16ads_v5":     4.192000,
		"Standard_E64-16as_v4":      4.032000,
		"Standard_E64-16as_v5":      3.616000,
		"Standard_E64-16ds_v4":      4.608000,
		"Standard_E64-16ds_v5":      4.608000,
		"Standard_E64-16s_v3":       3.629000,
		"Standard_E64-16s_v4":       4.032000,
		"Standard_E64-16s_v5":       4.032000,
		"Standard_E64-32ads_v5":     4.192000,
		"Standard_E64-32as_v4":      4.032000,
		"Standard_E64-32as_v5":      3.616000,
		"Standard_E64-32ds_v4":      4.608000,
		"Standard_E64-32ds_v5":      4.608000,
		"Standard_E64-32s_v3":       3.629000,
		"Standard_E64-32s_v4":       4.032000,
		"Standard_E64-32s_v5":       4.032000,
		"Standard_E64_v3":           3.629000,
		"Standard_E64_v4":           4.032000,
		"Standard_E64_v5":           4.032000,
		"Standard_E64a_v4":          4.032000,
		"Standard_E64ads_v5":        4.192000,
		"Standard_E64as_v4":         4.032000,
		"Standard_E64as_v5":         3.616000,
		"Standard_E64bds_v5":        5.344000,
		"Standard_E64bs_v5":         4.768000,
		"Standard_E64d_v4":          4.608000,
		"Standard_E64d_v5":          4.608000,
		"Standard_E64ds_v4":         4.608000,
		"Standard_E64ds_v5":         4.608000,
		"Standard_E64i_v3":          3.629000,
		"Standard_E64is_v3":         3.629000,
		"Standard_E64s_v3":          3.629000,
		"Standard_E64s_v4":          4.032000,
		"Standard_E64s_v5":          4.032000,
		"Standard_E8-2ads_v5":       0.524000,
		"Standard_E8-2as_v4":        0.504000,
		"Standard_E8-2as_v5":        0.452000,
		"Standard_E8-2ds_v4":        0.576000,
		"Standard_E8-2ds_v5":        0.576000,
		"Standard_E8-2s_v3":         0.504000,
		"Standard_E8-2s_v4":         0.504000,
		"Standard_E8-2s_v5":         0.504000,
		"Standard_E8-4ads_v5":       0.524000,
		"Standard_E8-4as_v4":        0.504000,
		"Standard_E8-4as_v5":        0.452000,
		"Standard_E8-4ds_v4":        0.576000,
		"Standard_E8-4ds_v5":        0.576000,
		"Standard_E8-4s_v3":         0.504000,
		"Standard_E8-4s_v4":         0.504000,
		"Standard_E8-4s_v5":         0.504000,
		"Standard_E80ids_v4":        5.760000,
		"Standard_E80is_v4":         5.040000,
		"Standard_E8_v3":            0.504000,
		"Standard_E8_v4":            0.504000,
		"Standard_E8_v5":            0.504000,
		"Standard_E8a_v4":           0.504000,
		"Standard_E8ads_v5":         0.524000,
		"Standard_E8as_v4":          0.504000,
		"Standard_E8as_v5":          0.452000,
		"Standard_E8bds_v5":         0.668000,
		"Standard_E8bs_v5":          0.596000,
		"Standard_E8d_v4":           0.576000,
		"Standard_E8d_v5":           0.576000,
		"Standard_E8ds_v4":          0.576000,
		"Standard_E8ds_v5":          0.576000,
		"Standard_E8pds_v5":         0.461000,
		"Standard_E8ps_v5":          0.403000,
		"Standard_E8s_v3":           0.504000,
		"Standard_E8s_v4":           0.504000,
		"Standard_E8s_v5":           0.504000,
		"Standard_E96-24ads_v5":     6.288000,
		"Standard_E96-24as_v4":      6.048000,
		"Standard_E96-24as_v5":      5.424000,
		"Standard_E96-24ds_v5":      6.912000,
		"Standard_E96-24s_v5":       6.048000,
		"Standard_E96-48ads_v5":     6.288000,
		"Standard_E96-48as_v4":      6.048000,
		"Standard_E96-48as_v5":      5.424000,
		"Standard_E96-48ds_v5":      6.912000,
		"Standard_E96-48s_v5":       6.048000,
		"Standard_E96_v5":           6.048000,
		"Standard_E96a_v4":          6.048000,
		"Standard_E96ads_v5":        6.288000,
		"Standard_E96as_v4":         6.048000,
		"Standard_E96as_v5":         5.424000,
		"Standard_E96bds_v5":        8.016000,
		"Standard_E96bs_v5":         7.152000,
		"Standard_E96d_v5":          6.912000,
		"Standard_E96ds_v5":         6.912000,
		"Standard_E96iads_v5":       6.917000,
		"Standard_E96ias_v5":        5.966000,
		"Standard_E96s_v5":          6.048000,
		"Standard_EC16ads_v5":       1.153000,
		"Standard_EC16as_v5":        0.994000,
		"Standard_EC20ads_v5":       1.441000,
		"Standard_EC20as_v5":        1.243000,
		"Standard_EC2ads_v5":        0.144000,
		"Standard_EC2as_v5":         0.124000,
		"Standard_EC32ads_v5":       2.306000,
		"Standard_EC32as_v5":        1.989000,
		"Standard_EC48ads_v5":       3.458000,
		"Standard_EC48as_v5":        2.983000,
		"Standard_EC4ads_v5":        0.288000,
		"Standard_EC4as_v5":         0.249000,
		"Standard_EC64ads_v5":       4.611000,
		"Standard_EC64as_v5":        3.978000,
		"Standard_EC8ads_v5":        0.576000,
		"Standard_EC8as_v5":         0.497000,
		"Standard_EC96ads_v5":       6.917000,
		"Standard_EC96as_v5":        5.966000,
		"Standard_EC96iads_v5":      7.608000,
		"Standard_EC96ias_v5":       6.563000,
		"Standard_F1":               0.049700,
		"Standard_F16":              0.796000,
		"Standard_F16s":             0.796000,
		"Standard_F16s_v2":          0.677000,
		"Standard_F1s":              0.049700,
		"Standard_F2":               0.099000,
		"Standard_F2s":              0.099000,
		"Standard_F2s_v2":           0.084600,
		"Standard_F32s_v2":          1.353000,
		"Standard_F4":               0.199000,
		"Standard_F48s_v2":          2.030000,
		"Standard_F4s":              0.199000,
		"Standard_F4s_v2":           0.169000,
		"Standard_F64s_v2":          2.706000,
		"Standard_F72s_v2":          3.045000,
		"Standard_F8":               0.398000,
		"Standard_F8s":              0.398000,
		"Standard_F8s_v2":           0.338000,
		"Standard_FX12mds":          1.116000,
		"Standard_FX24mds":          2.232000,
		"Standard_FX36mds":          3.348000,
		"Standard_FX48mds":          4.464000,
		"Standard_FX4mds":           0.372000,
		"Standard_H16":              1.807000,
		"Standard_H16m":             2.422000,
		"Standard_H16mr":            2.664000,
		"Standard_H16r":             1.988000,
		"Standard_H8":               0.904000,
		"Standard_H8m":              1.211000,
		"Standard_HB120-16rs_v2":    3.600000,
		"Standard_HB120-16rs_v3":    3.600000,
		"Standard_HB120-32rs_v2":    3.600000,
		"Standard_HB120-32rs_v3":    3.600000,
		"Standard_HB120-64rs_v2":    3.600000,
		"Standard_HB120-64rs_v3":    3.600000,
		"Standard_HB120-96rs_v2":    3.600000,
		"Standard_HB120-96rs_v3":    3.600000,
		"Standard_HB120rs_v2":       3.600000,
		"Standard_HB120rs_v3":       3.600000,
		"Standard_HB176-144rs_v4":   7.200000,
		"Standard_HB176-24rs_v4":    7.200000,
		"Standard_HB176-48rs_v4":    7.200000,
		"Standard_HB176-96rs_v4":    7.200000,
		"Standard_HB176rs_v4":       7.200000,
		"Standard_HB60-15rs":        2.280000,
		"Standard_HB60-30rs":        2.280000,
		"Standard_HB60-45rs":        2.280000,
		"Standard_HB60rs":           2.280000,
		"Standard_HC44-16rs":        3.168000,
		"Standard_HC44-32rs":        3.168000,
		"Standard_HC44rs":           3.168000,
		"Standard_HX176-144rs":      8.640000,
		"Standard_HX176-24rs":       8.640000,
		"Standard_HX176-48rs":       8.640000,
		"Standard_HX176-96rs":       8.640000,
		"Standard_HX176rs":          8.640000,
		"Standard_L16as_v3":         1.248000,
		"Standard_L16s_v2":          1.248000,
		"Standard_L16s_v3":          1.392000,
		"Standard_L32as_v3":         2.496000,
		"Standard_L32s_v2":          2.496000,
		"Standard_L32s_v3":          2.784000,
		"Standard_L48as_v3":         3.744000,
		"Standard_L48s_v2":          3.744000,
		"Standard_L48s_v3":          4.176000,
		"Standard_L64as_v3":         4.992000,
		"Standard_L64s_v2":          4.992000,
		"Standard_L64s_v3":          5.568000,
		"Standard_L80as_v3":         6.240000,
		"Standard_L80s_v2":          6.240000,
		"Standard_L80s_v3":          6.960000,
		"Standard_L88is_v2":         6.864000,
		"Standard_L8as_v3":          0.624000,
		"Standard_L8s_v2":           0.624000,
		"Standard_L8s_v3":           0.696000,
		"Standard_M128":             13.338000,
		"Standard_M128-32ms":        26.688000,
		"Standard_M128-64ms":        26.688000,
		"Standard_M128dms_v2":       26.690000,
		"Standard_M128ds_v2":        13.340000,
		"Standard_M128m":            26.688000,
		"Standard_M128ms":           26.688000,
		"Standard_M128ms_v2":        26.269000,
		"Standard_M128s":            13.338000,
		"Standard_M128s_v2":         12.919000,
		"Standard_M16-4ms":          3.073000,
		"Standard_M16-8ms":          3.073000,
		"Standard_M16ms":            3.073000,
		"Standard_M16s":             2.387000,
		"Standard_M192idms_v2":      32.064000,
		"Standard_M192ids_v2":       16.032000,
		"Standard_M192ims_v2":       31.643000,
		"Standard_M192is_v2":        15.611000,
		"Standard_M208-104ms_v2":    44.620000,
		"Standard_M208-104s_v2":     22.310000,
		"Standard_M208-52ms_v2":     44.620000,
		"Standard_M208-52s_v2":      22.310000,
		"Standard_M208ms_v2":        44.620000,
		"Standard_M208s_v2":         22.310000,
		"Standard_M32-16ms":         6.146000,
		"Standard_M32-8ms":          6.146000,
		"Standard_M32dms_v2":        6.146000,
		"Standard_M32ls":            2.873000,
		"Standard_M32ms":            6.146000,
		"Standard_M32ms_v2":         6.041000,
		"Standard_M32s":             3.335000,
		"Standard_M32ts":            2.707000,
		"Standard_M416-104ms_v2":    99.150000,
		"Standard_M416-104s_v2":     49.580000,
		"Standard_M416-208ms_v2":    99.150000,
		"Standard_M416-208s_v2":     49.580000,
		"Standard_M416is_v2":        49.580000,
		"Standard_M416ms_v2":        99.150000,
		"Standard_M416s_v2":         49.580000,
		"Standard_M64":              6.669000,
		"Standard_M64-16ms":         10.337000,
		"Standard_M64-32ms":         10.337000,
		"Standard_M64dms_v2":        10.340000,
		"Standard_M64ds_v2":         6.669000,
		"Standard_M64ls":            5.415000,
		"Standard_M64m":             10.337000,
		"Standard_M64ms":            10.337000,
		"Standard_M64ms_v2":         10.130000,
		"Standard_M64s":             6.669000,
		"Standard_M64s_v2":          6.459000,
		"Standard_M8-2ms":           1.536500,
		"Standard_M8-4ms":           1.536500,
		"Standard_M8ms":             1.536500,
		"Standard_NC12":             1.800000,
		"Standard_NC12s_v2":         4.140000,
		"Standard_NC12s_v3":         6.120000,
		"Standard_NC16as_T4_v3":     1.204000,
		"Standard_NC24":             3.600000,
		"Standard_NC24ads_A100_v4":  3.673000,
		"Standard_NC24r":            3.960000,
		"Standard_NC24rs_v2":        9.108000,
		"Standard_NC24rs_v3":        13.460000,
		"Standard_NC24s_v2":         8.280000,
		"Standard_NC24s_v3":         12.240000,
		"Standard_NC48ads_A100_v4":  7.346000,
		"Standard_NC4as_T4_v3":      0.526000,
		"Standard_NC6":              0.900000,
		"Standard_NC64as_T4_v3":     4.352000,
		"Standard_NC6s_v2":          2.070000,
		"Standard_NC6s_v3":          3.060000,
		"Standard_NC8as_T4_v3":      0.752000,
		"Standard_NC96ads_A100_v4":  14.692000,
		"Standard_ND12s":            4.140000,
		"Standard_ND24rs":           9.108000,
		"Standard_ND24s":            8.280000,
		"Standard_ND40rs_v2":        22.032000,
		"Standard_ND40s_v2":         12.240000,
		"Standard_ND6s":             2.070000,
		"Standard_ND96amsr_A100_v4": 32.770000,
		"Standard_ND96asr_A100_v4":  27.197000,
		"Standard_ND96asr_v4":       27.197000,
		"Standard_NP10s":            1.650000,
		"Standard_NP20s":            3.300000,
		"Standard_NP40s":            6.600000,
		"Standard_NV12":             2.280000,
		"Standard_NV12ads_A10_v5":   0.908000,
		"Standard_NV12s_v3":         1.140000,
		"Standard_NV16as_v4":        0.932000,
		"Standard_NV18ads_A10_v5":   1.600000,
		"Standard_NV24":             4.560000,
		"Standard_NV24s_v3":         2.280000,
		"Standard_NV32as_v4":        1.864000,
		"Standard_NV36adms_A10_v5":  4.520000,
		"Standard_NV36ads_A10_v5":   3.200000,
		"Standard_NV48s_v3":         4.560000,
		"Standard_NV4as_v4":         0.233000,
		"Standard_NV6":              1.140000,
		"Standard_NV6ads_A10_v5":    0.454000,
		"Standard_NV72ads_A10_v5":   6.520000,
		"Standard_NV8as_v4":         0.466000,
		"Standard_PB12s":            1.660000,
		"Standard_PB24s":            3.320000,
		"Standard_PB6s":             0.830000,
	}
}
