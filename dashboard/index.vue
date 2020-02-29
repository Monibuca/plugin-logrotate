<template>
    <Tabs  type="card">
        <TabPane label="日志文件" class="tabpanel">
            <List>
                <ListItem v-for="item in logFiles" :key="item.Name">
                    <ListItemMeta :title="item.Name" :description="networkFormat(item.Size)"></ListItemMeta>
                </ListItem>
            </List>
        </TabPane>
        <TabPane label="日志跟踪" class="tabpanel">
            <div>
                自动滚动
                <i-switch v-model="autoScroll" />
            </div>
            <div ref="logContainer" class="log-container">
                <pre><template v-for="item in logs">{{item+"\n"}}</template></pre>
            </div>
        </TabPane>
        <TabPane label="日志查询" class="tabpanel">
            <i-input search @on-search="onSearch"></i-input>
            <pre>{{result}}</pre>
        </TabPane>
    </Tabs>
</template>

<script>
let logsES = null;
export default {
    data() {
        return {
            autoScroll: true,
            logs: [],
            logFiles: [],
            result: ""
        };
    },
    mounted() {
        logsES = new EventSource("/logrotate/tail");
        logsES.onmessage = evt => {
            if (!evt.data) return;
            this.logs.push(evt.data);
        };
        window.ajax.getJSON("/logrotate/list").then(x => (this.logFiles = x));
    },
    deactivated() {
        logsES.close();
    },
    methods: {
        networkFormat(value, unit = "") {
            if (value > 1024 && uintInc[unit]) {
                return this.networkFormat(value / 1024, uintInc[unit]);
            }
            return value.toFixed(2).replace(".00", "") + unit + "B";
        },
        onSearch(value) {
            window.ajax
                .get("/logrotate/find?query=" + value)
                .then(x => (this.result = x));
        }
    },
    updated() {
        if (this.autoScroll) {
            this.$refs.logContainer.scrollTop = this.$refs.logContainer.offsetHeight;
        }
    }
};
</script>

<style>
@import url("/iview.css");
.log-container {
    overflow-y: auto;
    max-height: 500px;
}
.tabpanel {
    padding: 0 20px;
}
</style>