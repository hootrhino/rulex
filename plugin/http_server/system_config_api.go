package httpserver

import "github.com/gin-gonic/gin"

/*
*
* WIFI
*
 */
func SetWifi(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
	}

}

/*
*
* 设置时间、时区
*
 */
func SetTime(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
	}

}

/*
*
* 设置静态网络IP等

	{
	  "network": {
	    "version": 2,
	    "renderer": "networkd",
	    "ethernets": {
	      "enp0s9": {
	        "dhcp4": "no",
	        "addresses": [
	          "192.168.121.221/24"
	        ],
	        "gateway4": "192.168.121.1",
	        "nameservers": {
	          "addresses": [
	            "8.8.8.8",
	            "1.1.1.1"
	          ]
	        }
	      }
	    }
	  }
	}
*/
func SetStaticNetwork(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
	}

}
