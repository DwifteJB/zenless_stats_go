let wpRequire = webpackChunkdiscord_developers.push([[Symbol()], {}, r => r]);
webpackChunkdiscord_developers.pop();
let api = Object.values(wpRequire.c).find(x => x?.exports?.Bo?.get).exports.Bo;

const appId = "1524599002165674096";

console.log("[Widget Fix] Looking up existing widget config...");
const listRes = await api.get({url: `/applications/${appId}/widget-configs`});
const list = Array.isArray(listRes.body) ? listRes.body : (listRes.body.configs || listRes.body.items || [listRes.body]);
const configId = list[0].config_id || list[0].id;
console.log("[Widget Fix] Config id:", configId);

const dataField = name => ({presentation_type: "text", value_type: "data", value: name});
const stat = (value, label) => ({fields: {value: dataField(value), label: dataField(label)}});
const surfaces = {
	widget_top: {
		layout: "widget_top_hero",
		components: {
			hero_image: {fields: {image: {presentation_type: "image", value_type: "custom_string", value: "change this to an image"}}},
			title: {fields: {text: dataField("nickname")}}
		}
	},
	widget_bottom: {
		layout: "widget_bottom_stats",
		components: {
			stat_1: stat("IL", "IL_str"),
			stat_2: stat("ach", "ach_str"),
			stat_3: stat("SBT", "SBT_str"),
			stat_4: stat("ENER", "ENER_str"),
			stat_5: stat("days", "days_str"),
			stat_6: stat("polychromes", "polychromes_str")
		}
	},
	add_widget_preview: {
		layout: "add_widget_preview_hero",
		components: {
			hero_image: {fields: {image: {presentation_type: "image", value_type: "custom_string", value: "change this to an image"}}}
		}
	}
};

console.log("[Widget Fix] Binding fields to dynamic data...");
await api.patch({url: `/applications/${appId}/widget-configs/${configId}`, body: {surfaces}});
await api.post({url: `/applications/${appId}/widget-configs/${configId}/publish`});
console.log("[Widget Fix] Done. Run the Go app, then reopen your profile to see the stats.");
