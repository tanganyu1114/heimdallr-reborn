import { findHmdrHost } from '@/api/hmdr_host'
import { getOptions } from '@/api/hmdr_conf'

export default {
  methods: {
    async buildEnvDetailsList(detailRes) {
      if (!detailRes) return []

      const envDetailsObj = detailRes.env_details || {}
      const envDetailsList = []

      const hostIdSet = new Set()
      const hostIdToServerDetailMap = new Map()
      const serverMatchResults = new Map()

      for (const envId in envDetailsObj) {
        const envDetail = envDetailsObj[envId]
        if (envDetail.server_details && Array.isArray(envDetail.server_details)) {
          for (const serverDetail of envDetail.server_details) {
            if (serverDetail.server && serverDetail.server.host_id) {
              const hostId = serverDetail.server.host_id
              hostIdSet.add(hostId)
              hostIdToServerDetailMap.set(hostId, serverDetail)
            }
          }
        }
      }

      const serverIpMap = new Map()
      const serverPathMap = new Map()

      if (hostIdSet.size > 0) {
        const optionsRes = await getOptions().catch(err => {
          console.error('Failed to get options:', err)
          return null
        })

        if (optionsRes && optionsRes.data && Array.isArray(optionsRes.data)) {
          const bifrostPatterns = []

          for (const group of optionsRes.data) {
            const groupLabel = group.label || '-'
            if (group.children && Array.isArray(group.children)) {
              for (const host of group.children) {
                const hostLabel = host.label || '-'
                const hostValue = host.value
                if (host.children && Array.isArray(host.children)) {
                  for (const bifrost of host.children) {
                    const bifrostLabel = bifrost.label || '-'
                    const bifrostValue = bifrost.value
                    bifrostPatterns.push({
                      pattern: bifrostValue,
                      path: groupLabel + ' / ' + hostLabel + ' / ' + bifrostLabel,
                      hostValue: String(hostValue)
                    })
                  }
                }
              }
            }
          }

          for (const envId in envDetailsObj) {
            const envDetail = envDetailsObj[envId]
            if (envDetail.server_details && Array.isArray(envDetail.server_details)) {
              for (const serverDetail of envDetail.server_details) {
                if (!serverDetail.server || !serverDetail.server.host_id) continue

                const hostId = serverDetail.server.host_id
                const server = serverDetail.server
                const matchResult = this.findMatchingPath(server, bifrostPatterns, hostId)

                if (matchResult) {
                  const key = hostId + ':' + server.name
                  serverMatchResults.set(key, matchResult)
                  serverPathMap.set(key, matchResult.path)
                }
              }
            }
          }
        }

        const hostPromises = Array.from(hostIdSet).map(hId =>
          findHmdrHost({ ID: hId }).catch(err => {
            console.error('Failed to get host info for host_id ' + hId + ':', err)
            return null
          })
        )

        const hostResults = await Promise.all(hostPromises)
        for (const hostRes of hostResults) {
          if (hostRes && hostRes.data && hostRes.data.rehmdrHost) {
            const host = hostRes.data.rehmdrHost
            serverIpMap.set(host.ID, host.ipaddr || '-')
          }
        }
      }

      for (const envId in envDetailsObj) {
        const envDetail = envDetailsObj[envId]
        const serverDetailsArray = []

        if (envDetail.server_details && Array.isArray(envDetail.server_details)) {
          for (const serverDetail of envDetail.server_details) {
            const hostId = serverDetail.server && serverDetail.server.host_id
            const bifrostName = serverDetail.server ? (serverDetail.server.name || '-') : '-'
            const key = hostId + ':' + bifrostName

            if (!serverPathMap.has(key)) {
              continue
            }

            const serverIP = hostId ? (serverIpMap.get(hostId) || 'HostID: ' + hostId) : '-'
            const serverPath = serverPathMap.get(key)

            const proxiesObj = serverDetail.proxies || {}
            const proxiesArray = typeof proxiesObj === 'object' && !Array.isArray(proxiesObj)
              ? Object.values(proxiesObj).sort((a, b) => (a.proxy_index || 0) - (b.proxy_index || 0))
              : (Array.isArray(proxiesObj) ? proxiesObj : [])

            serverDetailsArray.push({
              server: serverDetail.server,
              server_ip: serverIP,
              server_path: serverPath,
              proxies: proxiesArray
            })
          }
        }

        if (serverDetailsArray.length > 0) {
          envDetailsList.push({
            env_id: parseInt(envId),
            server_details: serverDetailsArray
          })
        }
      }

      envDetailsList.sort((a, b) => a.env_id - b.env_id)
      return envDetailsList
    },
    findMatchingPath(server, bifrostPatterns, hostId) {
      const regexPatterns = []

      if (server.http_proxy_server_name_match_str) {
        regexPatterns.push(server.http_proxy_server_name_match_str)
      }

      if (server.stream_proxy_server_name_match_str) {
        regexPatterns.push(server.stream_proxy_server_name_match_str)
      }

      if (regexPatterns.length === 0) {
        return null
      }

      for (const pattern of bifrostPatterns) {
        if (String(hostId) !== pattern.hostValue) {
          continue
        }

        for (const regexStr of regexPatterns) {
          try {
            const regex = new RegExp(regexStr)
            if (regex.test(pattern.pattern)) {
              return { path: pattern.path, matchedField: pattern.pattern, regex: regexStr }
            }
          } catch (e) {
            if (pattern.pattern.includes(regexStr)) {
              return { path: pattern.path, matchedField: pattern.pattern, regex: regexStr }
            }
          }
        }
      }

      return null
    }
  }
}
