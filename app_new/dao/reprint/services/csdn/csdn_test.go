package csdn

import "testing"

func TestReq(t *testing.T) {
	title := "test_golang版"
	content := "<p>成功啦</p>\n"
	description := "发送成功（转载版）"
	tags := "go"
	cookie := "uuid_tt_dd=10_20879684060-1621692917600-112215; __gads=ID=c9a6da06f2cff7f2-227f85bc47c800ad:T=1621692923:S=ALNI_MYmpRXMIHe878uMQ_HUrxZB_WpDJQ; UserName=qq_39122387; UserInfo=a318ceb7325646d6813136b47e299dff; UserToken=a318ceb7325646d6813136b47e299dff; UserNick=屌丝_Asa; AU=36B; UN=qq_39122387; BT=1640346910339; p_uid=U010000; ssxmod_itna=eqmxuD9D0Qi=D=IK0Ly7tiNpghDBnDmuuyq9D0=7DlgQxA5D8D6DQeGTTnp3Te3Yni2vWjK8Duw=jeEixLWt7OYP4FiZ1xiTD4q07Db4GkDAqiOD7kRvomrvKD3Dm4i3DDoQDgDmKGg+qDf94NDG3=DfuXMDBYZS4DBmhXXFeGvGlj27eDIUqG+KqGuBkju7g=cRkjXD0poQDf7v0548LxHQ7GNYhq=/oPE8iwKnZv53/wt/GD40BgdWA0x7B5QeD===; ssxmod_itna2=eqmxuD9D0Qi=D=IK0Ly7tiNpghDBnDmuuyq9D0=D66tlEx0vF303bsce3RD6QALIf7ZzWqhiee=CT73+e4tlvivdpfWrcEQES7iQ6zuFukw0px/+Yi6qcSIl7ApwR6UczrvHQ60NepwzorWq/pxk+EKGZGes+WwkQRWCxKohZiOzYg9TWaOIH+BvgxmHNxvaE+A+S7D88rDGLtMpNY4j41FLLrfqgU9f3rpHcS9O/+Shz9Kh4wKa1Qw=SpbtR82qixex9lkN=De4itWueDbQMOvbfTArqtnM9ZnzfZmEfEWHLv/RZs8Q8K+9va85+z5TK+xsbHQ4+3i8IyQ3whKDUKD/tFD=nNCAbAFPXHHKipCmqxecCBTm5hDohYB8KEU+2MWgVm7HbBZlvEbAbQPVQcoq7OnonCBOS4m45rhvCjHThSuBpRSe0z8E0uA0ysFxDKuDlDQmjK=Ax8r4KojY7jD2jA6xQgqeDDLxD2WGDD==; Hm_up_6bcd52f51e9b3dce32bec4a3997715ac={\"islogin\":{\"value\":\"1\",\"scope\":1},\"isonline\":{\"value\":\"1\",\"scope\":1},\"isvip\":{\"value\":\"0\",\"scope\":1},\"uid_\":{\"value\":\"qq_39122387\",\"scope\":1}}; Hm_ct_6bcd52f51e9b3dce32bec4a3997715ac=6525*1*10_20879684060-1621692917600-112215!5744*1*qq_39122387; dc_session_id=10_1648470227770.123525; c_first_ref=default; c_first_page=https://www.csdn.net/; c_segment=7; c_page_id=default; dc_sid=506efe62be0d01799187f9d16a27f80d; Hm_lvt_6bcd52f51e9b3dce32bec4a3997715ac=1647696313,1648470227; log_Id_view=451; c_pref=https://im.csdn.net/chat/csdn_sysnotify; c_ref=https://mp.csdn.net/mp_blog/creation/editor; dc_tos=r9gfz5; log_Id_pv=121; Hm_lpvt_6bcd52f51e9b3dce32bec4a3997715ac=1648470498; log_Id_click=41"

	param, _ := NewParam(title,content, tags, Description(description))

	msg, err := PublishArticle(param, cookie)

	if err != nil {
		t.Errorf("publish article error, it is %v", err)
		return
	}

	if msg != "发布成功。" {
		t.Errorf("error, massage is %v", msg)
	}

}