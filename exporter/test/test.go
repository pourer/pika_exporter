package test

var InfoCases = []struct {
	Name string
	Info string
}{
	{"v2.2.6_master", V226MasterInfo},
	{"v2.2.6_slave", V226SlaveInfo},

	{"v2.3.6_master", V236MasterInfo},
	{"v2.3.6_slave", V236SlaveInfo},

	{"v3.0.10_master", V3010MasterInfo},
	{"v3.0.10_slave", V3010SlaveInfo},

	{"v3.0.16_master", V3016MasterInfo},
	{"v3.0.16_slave", V3016SlaveInfo},

	{"v3.1.0_master", V310MasterInfo},
	{"v3.1.0_slave", V310SlaveInfo},

	{"v3.2.0_master", V320MasterInfo},
	{"v3.2.0_slave", V320SlaveInfo},

	{"v3.2.7_slave", V327SlaveInfo},

	{"v3.3.5_master", V335MasterInfo},
	{"v3.3.5_slave", V335SlaveInfo},
}
