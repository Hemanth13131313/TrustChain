package trustchain

import data.trustchain.rules

# Default deny
default allow = false
default reason = "no policy evaluated to true"

# The entrypoint rule
allow {
    not has_critical_vulnerabilities
    has_valid_signature
    meets_slsa_requirements
}

has_critical_vulnerabilities {
    # simulated vulnerability data passed in input context
    v := input.vulnerabilities[_]
    v.severity == "CRITICAL"
    v.status == "UNMITIGATED"
}

has_valid_signature {
    input.signatures_count > 0
    input.signature_verified == true
}

meets_slsa_requirements {
    input.slsa_level >= 3
}

reason = "critical unmitigated vulnerabilities found" {
    has_critical_vulnerabilities
} else = "no valid signatures" {
    not has_valid_signature
} else = "insufficient SLSA level (requires >= 3)" {
    not meets_slsa_requirements
} else = "allowed" {
    allow
}
