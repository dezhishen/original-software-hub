export function detectCurrentPlatform() {
  const ua = String(navigator.userAgent || '').toLowerCase()
  const platform = String(
    (navigator.userAgentData && navigator.userAgentData.platform) || navigator.platform || ''
  ).toLowerCase()
  const src = `${ua} ${platform}`
  if (/iphone|ipad|ipod/.test(src)) return { id: 'ios', label: 'iOS' }
  if (/android/.test(src)) return { id: 'android', label: 'Android' }
  if (/mac|darwin/.test(src)) return { id: 'macos', label: 'macOS' }
  if (/win/.test(src)) return { id: 'windows', label: 'Windows' }
  if (/linux|x11/.test(src)) return { id: 'linux', label: 'Linux' }
  return { id: 'web', label: 'Web' }
}

export function detectCurrentArchitecture() {
  const ua = String(navigator.userAgent || '').toLowerCase()
  const uaArch = String(
    (navigator.userAgentData && navigator.userAgentData.architecture) || ''
  ).toLowerCase()
  const platform = String(
    (navigator.userAgentData && navigator.userAgentData.platform) || navigator.platform || ''
  ).toLowerCase()
  const bitness = String(
    (navigator.userAgentData && navigator.userAgentData.bitness) || ''
  ).toLowerCase()
  const src = `${uaArch} ${platform} ${bitness} ${ua}`
  if (/arm64|aarch64/.test(src)) return { id: 'arm64', label: 'ARM64' }
  if (/x86_64|x64|amd64|wow64/.test(src) || bitness === '64') return { id: 'x64', label: 'x64' }
  if (/x86|i[3-6]86|ia32/.test(src)) return { id: 'x86', label: 'x86' }
  if (/universal/.test(src)) return { id: 'universal', label: 'Universal' }
  return { id: 'unknown', label: '未知' }
}

export function normalizePlatformId(label) {
  const v = String(label || '').trim().toLowerCase()
  if (v.includes('windows') || v === 'win') return 'windows'
  if (v.includes('macos') || v.includes('mac os') || v === 'mac' || v.includes('darwin')) return 'macos'
  if (v.includes('linux')) return 'linux'
  if (v.includes('android')) return 'android'
  if (v.includes('ios') || v.includes('iphone') || v.includes('ipad')) return 'ios'
  return v || 'other'
}

export function normalizeArchitectureId(label) {
  const v = String(label || '').trim().toLowerCase()
  if (/arm64|aarch64/.test(v)) return 'arm64'
  if (/x86_64|x64|amd64/.test(v)) return 'x64'
  if (/x86|i[3-6]86|ia32/.test(v)) return 'x86'
  if (/universal/.test(v)) return 'universal'
  return 'unknown'
}

export function architectureScore(variantArch, currentArchId) {
  const archId = normalizeArchitectureId(variantArch)
  if (archId === 'universal') return 85
  if (archId === currentArchId) return 100
  if (currentArchId === 'x64' && archId === 'x86') return 70
  if (currentArchId === 'arm64' && archId === 'x64') return 55
  if (currentArchId === 'x86' && archId === 'x64') return 60
  if (archId === 'unknown') return 45
  return 40
}

export function usePlatform() {
  return {
    currentPlatform: detectCurrentPlatform(),
    currentArchitecture: detectCurrentArchitecture()
  }
}
