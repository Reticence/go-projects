// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2017-11-07 15:57
// @Version : 0.1
// @Software: GoLand

package exercises

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type JsonData struct {
	Title       string   `json:"标题"`
	PicturePath []string `json:"图片路径"`
	Paragraphs  []string `json:"段落组"`
}

func save() {
	jjh := JsonData{"中国健康促进基金会健康管理研究所简介", []string{"T:/go-io-dir/jjh.png"}, make([]string, 4)}
	jjh.Paragraphs[0] = "为顺应现代医学目的和医学模式的转变，促进我国健康管理事业的发展，加强健康管理医学服务模式和适宜技术的研究，提高健康管理机构的服务品质和能力，推动健康管理学科体系建设、集成创新健康管理适宜技术和规范，开展多中心应用研究课题项目与转化，中国健康促进基金会于2012年正式设立了健康管理研究所，整合发挥专家、机构和企业的资源，构建“产-学-研-用”互惠多赢的机制与平台，坚定不移、百折不饶地开展健康管理示范基地建设“十百千”工程。"
	jjh.Paragraphs[1] = "学科建设为龙头，科研引领、人才支撑，“换思维、创条件、搭平台、建机制、谋发展”，开展全国健康管理学科建设科研示范课题项目是健康管理示范基地建设与发展的核心竞争力；我们开展应用研究课题目的不是单纯为科研而科研，而是重在转化应用，为转化申报国家课题打基础。学科发展的经验表明：实现医教研、适宜技术转化应用才是健康管理学科发展的催化剂和驱动力！"
	jjh.Paragraphs[2] = "健康管理研究所的核心任务：集成创新健康管理适宜技术与干预产品，组织开展相关多中心应用研究课题，紧紧围绕“人才、技术和规范管理”三个关键，推动健康管理理论与实践研究与成果转化，形成健康管理大数据。"
	jjh.Paragraphs[3] = "学术引领产业发展、产业促进学术进步！"

	meridian := JsonData{"经纶世纪医疗网络技术（北京）有限公司简介", []string{"T:/go-io-dir/meridian.png", "T:/go-io-dir/meridian_1.png"}, make([]string, 15)}
	meridian.Paragraphs[0] = "经纶世纪医疗网络技术（北京）有限公司（简称“经纶世纪”）由留美归国博士团队创建，致力于将物联网、大数据人工智能、智能服务机器人等高科技与健康医疗服务深度融合，为医院/体检中心、保险公司、企业、社区提供大数据与健康医疗人工智能辅助决策与管理系统。"
	meridian.Paragraphs[1] = "三位一体核心技术与产品："
	meridian.Paragraphs[2] = "经纶世纪自主研发“健康医疗大数据、医学AI引擎、智能医学服务系统”三位一体的健康医疗服务平台与解决方案，包括："
	meridian.Paragraphs[3] = "健康医疗大数据系统平台（大数据清洗与标准化、机器学习分析、疾病预测模型、健康评估关联图）；"
	meridian.Paragraphs[4] = "智能健康物联网系统（亚健康、高血压、糖尿病等慢病管理服务）；"
	meridian.Paragraphs[5] = "智能癌症早期筛查与康复管理系统（癌症诊疗指南与临床路径智能管理，如甲状腺结节风险评估、甲状腺癌术后康复管理等）；"
	meridian.Paragraphs[6] = "智能服务机器人（全科医生智能助手机器人—慢性病/常见病智能诊疗、智能分级诊疗、远程医疗；居家养老智能伴侣机器人—健康、医疗、关怀）。"
	meridian.Paragraphs[7] = "公司荣誉资质："
	meridian.Paragraphs[8] = "北京市企业“未来之星”"
	meridian.Paragraphs[9] = "中关村国家自主创新园区高端领军人才"
	meridian.Paragraphs[10] = "中国健康促进基金会健康体检大数据和健康物联网平台"
	meridian.Paragraphs[11] = "北京海外高层次人才协会英才榜"
	meridian.Paragraphs[12] = "中国卫生信息学会健康医疗大数据应用评估和保障专业委员会 常委"
	meridian.Paragraphs[13] = "中国电子学会健康物联网专委会 专家"
	meridian.Paragraphs[14] = "中国医疗大数据与人工智能优秀实践案例奖"

	datas := map[string]JsonData{"基金会简介": jjh, "经纶世纪简介": meridian}
	bytes, _ := json.Marshal(datas)
	ioutil.WriteFile("T:/go-io/jd.json", bytes, 0666)
}

func RunJson() {
	//save()
	var jd JsonData
	bytes, _ := ioutil.ReadFile("T:/go-io/jd.json")
	if err := json.Unmarshal(bytes, &jd); err != nil {
		panic(err)
	}

	fmt.Println(jd)
}
