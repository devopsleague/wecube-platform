//
// This file was generated by the JavaTM Architecture for XML Binding(JAXB) Reference Implementation, v2.2.8-b130911.1802 
// See <a href="http://java.sun.com/xml/jaxb">http://java.sun.com/xml/jaxb</a> 
// Any modifications to this file will be lost upon recompilation of the source schema. 
// Generated on: 2020.11.09 at 03:29:54 PM CST 
//


package com.webank.wecube.platform.core.service.plugin.xml.register;

import java.util.ArrayList;
import java.util.List;
import javax.xml.bind.annotation.XmlAccessType;
import javax.xml.bind.annotation.XmlAccessorType;
import javax.xml.bind.annotation.XmlAttribute;
import javax.xml.bind.annotation.XmlElement;
import javax.xml.bind.annotation.XmlType;


/**
 * <p>Java class for authorityType complex type.
 * 
 * <p>The following schema fragment specifies the expected content contained within this class.
 * 
 * <pre>
 * &lt;complexType name="authorityType">
 *   &lt;complexContent>
 *     &lt;restriction base="{http://www.w3.org/2001/XMLSchema}anyType">
 *       &lt;sequence>
 *         &lt;element name="menu" type="{}menuType" maxOccurs="unbounded"/>
 *       &lt;/sequence>
 *       &lt;attribute name="systemRoleName" type="{http://www.w3.org/2001/XMLSchema}string" />
 *     &lt;/restriction>
 *   &lt;/complexContent>
 * &lt;/complexType>
 * </pre>
 * 
 * 
 */
@XmlAccessorType(XmlAccessType.FIELD)
@XmlType(name = "authorityType", propOrder = {
    "menu"
})
public class AuthorityType {

    @XmlElement(required = true)
    protected List<MenuType> menu;
    @XmlAttribute(name = "systemRoleName")
    protected String systemRoleName;

    /**
     * Gets the value of the menu property.
     * 
     * <p>
     * This accessor method returns a reference to the live list,
     * not a snapshot. Therefore any modification you make to the
     * returned list will be present inside the JAXB object.
     * This is why there is not a <CODE>set</CODE> method for the menu property.
     * 
     * <p>
     * For example, to add a new item, do as follows:
     * <pre>
     *    getMenu().add(newItem);
     * </pre>
     * 
     * 
     * <p>
     * Objects of the following type(s) are allowed in the list
     * {@link MenuType }
     * 
     * 
     */
    public List<MenuType> getMenu() {
        if (menu == null) {
            menu = new ArrayList<MenuType>();
        }
        return this.menu;
    }

    /**
     * Gets the value of the systemRoleName property.
     * 
     * @return
     *     possible object is
     *     {@link String }
     *     
     */
    public String getSystemRoleName() {
        return systemRoleName;
    }

    /**
     * Sets the value of the systemRoleName property.
     * 
     * @param value
     *     allowed object is
     *     {@link String }
     *     
     */
    public void setSystemRoleName(String value) {
        this.systemRoleName = value;
    }

}