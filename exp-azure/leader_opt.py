import sys
import pandas as pd
from geopy.distance import geodesic

# 2つの地点の緯度経度を定義

# 結果を出力
region_name_table={
    "francecentral":"France Central",
    "southafricanorth":"South Africa North",
    "australiasoutheast":"Australia SouthEast",
    "brazilsouth":"Brazil South",
    "canadacentral":"Canada Central",
    "japaneast":"Japan East",
    "koreacentral":"Korea Central",
    "southeastasia":"South East Asia",
    "uksouth":"UK South",
    "eastus":"East US",
    "westus":"West US",
    "southindia":"South India",
} 

positions = {
    "francecentral": (46.2276, 2.2137),
    "southafricanorth": (-25.731340, 28.218370),
    "australiasoutheast": (-37.8136, 144.9631),
    "brazilsouth": (-23.550520, -46.633309),
    "canadacentral": (43.6532, -79.3832),
    "japaneast": (35.68, 139.77),
    "koreacentral": (37.5665, 126.9780),
    "southeastasia": (1.3521, 103.8198),
    "uksouth": (51.5074, -0.1278),
    "eastus": (33.7490, -84.3880),
    "westus": (37.7749, -122.4194),
    "southindia": (12.8536, 80.8312),
    "centralus": (41.8781, -87.6298),
    "swedencentral": (59.3293, 18.0686)
}


def Leader(R, c):
        T_req=[geodesic(positions[r], positions[c]) for r in R]
        index = T_req.index(min(T_req))
        return R[index]
    
def main():
    args = sys.argv[1:]
    if len(args) != 2:
        print("引数は2つ必要です。")
        return
    client = args[0]
    replica = args[1].split(",")[:-1]
    print(Leader(replica,client))

if __name__ == "__main__":
    main()
