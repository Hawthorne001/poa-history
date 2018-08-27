package handlers

import (
	"github.com/gin-gonic/gin"
	"../assets"
	"../tree/tree"
	"../tree/content"
	"../../ethercrypto/web3history"
	"net/http"
	"encoding/hex"
	"log"
	"encoding/json"
)

type Proof struct {
	//Number string
	Hash   string
}
type Proofs []Proof

//UpdateAssetId Add new asset to assetId and change merkle tree
func UpdateAssetId(c *gin.Context) {
	if assets.Check(c) {
		assets.UpdateAssetsByAssetId(c)
		tx := assets.GetTxNumber(c)
		tx++
		content.AddContent(c, tx)
		root := tree.GetRoot(c)
		web3history.SendNewRootHash(root)
		defer assets.IncrementAssetTx(c)
	} else {
	}
}

//CreateAssetId Create new assetId with asset
func CreateAssetId(c *gin.Context) {
	id, er, try := assets.CheckAndReturn(c)
	if try {
		tx := assets.GetTxNumber(c)
		content.AddContent(c, tx)
		root := tree.GetRoot(c)
		web3history.SendNewRootHash(root)
		if er == "err" {
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"assetId":    id[0],
			"txNumber":   tx,
			"hash":       id[1],
			"merkleRoot": hex.EncodeToString(root),
		})
		return
	}
}

//List Lists all assets in DB
func List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"assets": assets.FindALlAssets(c),
	})
}

//GetTotalProof Get total Merkle proof
func GetTotalProof(c *gin.Context) {
	d := tree.GetProofs(c)
	var proofs = Proofs{}
	for i := 0; i < len(d); i++ {
		proofs = append(proofs,
			Proof{
				hex.EncodeToString(d[i]),
			})
	}

	myJson, err := json.Marshal(proofs)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}

	c.Data(http.StatusOK, "JSON", myJson)

}

//GetData Get timestamp and hash of specified asset in assetId
func GetData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"assets": hex.EncodeToString(assets.GetAssetByAssetIdAndTxNumber(c)),
	})
}
