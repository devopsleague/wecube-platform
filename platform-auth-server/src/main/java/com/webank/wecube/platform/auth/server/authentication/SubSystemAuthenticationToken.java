package com.webank.wecube.platform.auth.server.authentication;

import java.util.Collection;

import org.springframework.security.authentication.AbstractAuthenticationToken;
import org.springframework.security.core.GrantedAuthority;

public class SubSystemAuthenticationToken extends AbstractAuthenticationToken {
    private static final long serialVersionUID = -6762939343698545817L;

    private final Object principal;

    private final Object credentials;
    
    private final String nonce;

    public SubSystemAuthenticationToken(String aPrincipal, String aCredentials, String aNonce,
            Collection<? extends GrantedAuthority> anAuthorities) {
        super(anAuthorities);
        this.principal = aPrincipal;
        this.credentials = aCredentials;
        this.nonce = aNonce;
        super.setAuthenticated(true);
    }

    public SubSystemAuthenticationToken(String aPrincipal, String aCredentials, String aNonce) {
        super(null);
        this.principal = aPrincipal;
        this.credentials = aCredentials;
        this.nonce = aNonce;
        setAuthenticated(false);
    }

    @Override
    public Object getCredentials() {
        return credentials;
    }

    @Override
    public Object getPrincipal() {
        return principal;
    }

    public String getNonce() {
        return nonce;
    }
}
