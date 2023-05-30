//go:build windows

package windows

func Windowshtml() string {
	return `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Windows等级保护核查结果</title>
    <style>
        body {
            display: grid;
            grid-template-columns: 1fr 200px;
            gap: 10px;
            font-family: Arial, sans-serif;
            position: relative;
        }

        table {
            border-collapse: collapse;
            margin-bottom: 20px;
            width: 100%;
            box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
        }

        th,
        td {
            border: 1px solid #ddd;
            padding: 15px;
            text-align: left;
        }

        th {
            background-color: #007BFF;
            color: white;
            font-weight: bold;
        }

        tr:nth-child(even) {
            background-color: #f9f9f9;
        }

        tr:hover {
            background-color: #e6f2ff;
        }

        .watermark {
            font-size: 36px;
            color: rgba(128, 128, 128, 0.2);
            position: absolute;
            z-index: -1;
            transform: rotate(-30deg);
        }

        #toc {
            position: fixed;
            top: 20px;
            right: 30px;
            padding-left: 10px;
            background-color: #f8f9fa;
            padding: 10px;
            border: 1px solid #dee2e6;
            border-radius: 5px;
            box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
        }

        #toc ul {
            list-style-type: none;
            padding: 0;
        }

        #toc a {
            text-decoration: none;
            color: #333;
            display: block;
        }

        #toc a:hover {
            color: #007BFF;
        }
		.watermark {
			font-size: 36px;
			color: rgba(128, 128, 128, 0.2);
			position: absolute;
			z-index: 1000;
			transform: rotate(-30deg);
    	}
        pre {
            background-color: #f8f9fa;
            border: 1px solid #dee2e6;
            padding: 15px;
            border-radius: 5px;
            overflow-x: auto;
            font-family: "Courier New", Courier, monospace;
			white-space: pre-wrap;
			word-break: break-word;
        }
    </style>

<body>

    <div id="content">
        <h2 id="osinfo">操作系统信息</h2>
        <table>
            <thead>
                <tr>
                    <th>操作系统名称</th>
                    <th>版本</th>
                    <th>架构</th>
                    <th>安装日期</th>
                </tr>
            </thead>
            <tbody>
                操作系统详细信息
            </tbody>
        </table>

        <h2 id="user">用户信息</h2>
        <table>
            <thead>
                <tr>
                    <th>用户</th>
                    <th>全名</th>
					<th>SID</th>
                    <th>注释</th>
                    <th>启用</th>
                    <th>帐户到期</th>
                    <th>上次修改密码时间</th>
                    <th>需要密码</th>
                    <th>密码到期</th>
                    <th>上次登录时间</th>
                    <th>本地组</th>
                </tr>
            </thead>
            <tbody>
                用户详细信息
            </tbody>
        </table>

        <h2 id="password-check">密码复杂度检查</h2>
        <table>
            <thead>
                <tr>
                    <th>检查项</th>
                    <th>检查结果</th>
                    <th>是否符合</th>
                    <th>建议结果</th>
                </tr>
            </thead>
            <tbody>
                密码复杂度结果
            </tbody>
        </table>

        <h2 id="mstsc">远程桌面</h2>
        <table>
            <thead>
                <tr>
                    <th>检查项</th>
                    <th>是否关闭</th>
                    <th>建议结果</th>
                </tr>
            </thead>
            <tbody>
				<tr>
					<td>是否开启远程桌面</td>
					<td>开启远程桌面结果</td>
					<td>根据业务场景判断是否有必要开启</td>
				</tr>
            </tbody>
        </table>

        <h2 id="password-accounts">密码有限期检查</h2>

        <table>
            <thead>
                <tr>
                    <th>检查项</th>
                    <th>检查结果</th>
                    <th>是否符合</th>
                    <th>建议结果</th>
                </tr>
            </thead>
            <tbody>
                密码有效期检查结果
            </tbody>
        </table>

        <h2 id="lockout-check">失败锁定次数</h2>

        <table>
            <thead>
                <tr>
                    <th>检查项</th>
                    <th>检查结果</th>
                    <th>是否符合</th>
                    <th>建议结果</th>
                </tr>
            </thead>
            <tbody>
                失败锁定结果
            </tbody>
        </table>

        <h2 id="auditd">审计相关核查</h2>
        <table>
            <thead>
                <tr>
                    <th>检查项</th>
                    <th>检查结果</th>
                    <th>是否符合</th>
                    <th>建议结果</th>
                </tr>
            </thead>
            <tbody>
                审计相关结果
            </tbody>
        </table>
        <h2 id="highauditd">高级审计策略</h2>
        <table>
            <thead>
                <tr>
                    <th>检查项</th>
                    <th>检查结果</th>
                    <th>是否符合</th>
                    <th>建议结果</th>
                </tr>
            </thead>
            <tbody>
                高级审计策略结果
            </tbody>
        </table>

		<h2 id="auditattribute">日志属性</h2>
        <pre><code>日志属性结果
字段解释：
name: 标识配置文件的名称
enabled: 表示此配置是否启用
type: 配置文件的类型
owningPublisher: 此项应包含发布者的信息
isolation: 描述配置的隔离层级
channelAccess: 定义了如何访问该通道的权限设置。各段括号内的权限是按照特定顺序和格式排列的，每一段表示不同的用户或组的权限。
logging: 包含日志记录的相关信息：
logFileName: 日志文件的存储路径
retention: 是否保留旧日志。
autoBackup: 是否自动备份日志文件
maxSize: 日志文件的最大大小（以字节为单位）
publishing: 包含发布相关设置的信息：
fileMax: 最大文件数量
		</code></pre>


        <h2 id="screen">屏幕保护核查</h2>
        <table>
            <thead>
                <tr>
                    <th>检查项</th>
                    <th>检查结果</th>
                    <th>是否符合</th>
                    <th>建议结果</th>
                </tr>
            </thead>
            <tbody>
                屏幕保护相关结果
            </tbody>
        </table>
        <h2 id="iptables">防火墙状态检查</h2>
        <table>
            <thead>
                <tr>
                    <th>检查项</th>
                    <th>检查结果</th>
                    <th>是否符合</th>
                    <th>建议结果</th>
                </tr>
            </thead>
            <tbody>
                防火墙状态检查结果
            </tbody>
        </table>

		<h2 id="group">群组信息</h2>
        <pre><code>群组信息结果</code></pre>

		<h2 id="computer">防病毒</h2>
        <pre><code>防病毒结果
字段解释：
AMServiceEnabled: 杀毒软件服务是否启用。
AntivirusEnabled: 是否启用了杀毒软件。
AntivirusSignatureLastUpdated: 杀毒软件定义上次更新的时间。
AntispywareEnabled: 反间谍软件是否启用。
BehaviorMonitorEnabled: 行为监视是否启用。
FullScanAge: 上次全盘扫描的天数。
FullScanEndTime: 全盘扫描的结束时间。
FullScanStartTime: 全盘扫描的开始时间。
IoavProtectionEnabled: IOAV保护是否启用。
NISEnabled: 网络检查系统（NIS）是否启动。
NISEngineVersion: NIS引擎版本。
NISSignatureAge: NIS签名的天数。
NISSignatureLastUpdated: NIS签名上次更新的时间。
NISSignatureVersion: NIS签名版本。
OnAccessProtectionEnabled: 实时保护是否启用。
QuickScanAge: 上次快速扫描的天数。
QuickScanEndTime: 快速扫描的结束时间。
QuickScanStartTime: 快速扫描的开始时间。
		</code></pre>

		<h2 id="netshare">共享资源</h2>
        <pre><code>共享资源结果</code></pre>

		<h2 id="network">联网测试</h2>
        <pre><code>联网测试结果</code></pre>

		<h2 id="OptionalFeature">安装组件</h2>
        <pre><code>安装组件结果</code></pre>

		<h2 id="installer">安装程序</h2>
        <pre><code>安装程序结果</code></pre>

		<h2 id="systeminfo">系统信息</h2>
        <pre><code>系统信息结果</code></pre>

		<h2 id="tasklist">进程列表</h2>
        <pre><code>进程列表结果</code></pre>

        <h2 id="port">开放端口</h2>
        <pre><code>端口相关结果</code></pre>

		<h2 id="Service">Service</h2>
        <pre><code>Service结果</code></pre>

		<h2 id="schtasks">定时任务</h2>
        <pre><code>定时任务结果</code></pre>

		<h2 id="patch">已安装补丁信息</h2>
        <pre><code>补丁相关结果</code></pre>

        <h2 id="domainrlue">核查域防火墙规则</h2>
        <pre><code>域防火墙规则结果</code></pre>
        <h2 id="privaterlue">核查专网防火墙规则</h2>
        <pre><code>专网防火墙规则结果</code></pre>
        <h2 id="publicrlue">核查公共防火墙规则</h2>
        <pre><code>公共防火墙规则结果</code></pre>


    </div>

    <div id="toc">
        <h3>目录</h3>
        <ul>
            <li><a href="#osinfo">操作系统信息</a></li>
            <li><a href="#user">用户信息</a></li>
            <li><a href="#group">群组信息</a></li>
            <li><a href="#password-accounts">密码有效期</a></li>
            <li><a href="#password-check">密码复杂度</a></li>
            <li><a href="#lockout-check">失败锁定</a></li>
            <li><a href="#mstsc">远程桌面</a></li>
            <li><a href="#auditd">审计策略</a></li>
            <li><a href="#highauditd">高级审计策略</a></li>
            <li><a href="#auditattribute">日志属性</a></li>
            <li><a href="#screen">屏幕保护</a></li>
            <li><a href="#port">开放端口</a></li>
            <li><a href="#computer">防病毒</a></li>
            <li><a href="#netshare">共享资源</a></li>
            <li><a href="#systeminfo">系统信息</a></li>
            <li><a href="#tasklist">进程列表</a></li>
            <li><a href="#installer">已安装程序</a></li>
            <li><a href="#OptionalFeature">已安装组件</a></li>
            <li><a href="#Service">Service</a></li>
            <li><a href="#schtasks">定时任务</a></li>
            <li><a href="#patch">安装补丁信息</a></li>
            <li><a href="#network">联网测试</a></li>
            <li><a href="#iptables">防火墙状态</a></li>
            <li><a href="#domainrlue">域防火墙规则</a></li>
            <li><a href="#privaterlue">专网防火墙规则</a></li>
            <li><a href="#publicrlue">公共防火墙规则</a></li>
        </ul>
    </div>

	<div class="watermark" style="top: 10%; left: 15%;">高业尚-SelinuxG</div>
	<div class="watermark" style="top: 20%; left: 40%;">高业尚-SelinuxG</div>
	<div class="watermark" style="top: 50%; left: 10%;">高业尚-SelinuxG</div>
	<div class="watermark" style="top: 80%; left: 65%;">高业尚-SelinuxG</div>
	<div class="watermark" style="top: 30%; left: 20%;">高业尚-SelinuxG</div>


</body>

</html>
`
}
