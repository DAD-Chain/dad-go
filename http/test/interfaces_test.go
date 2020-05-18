package test

import (
	"testing"
	"bytes"
	"fmt"
	"os"
	"github.com/dad-go/core/types"
	"github.com/dad-go/common"
)

func TestTxDeserialize(t *testing.T) {
	bys, _ := common.HexToBytes("00d06de7c94900fdcf2b0138c56b6c766b00527ac46c766b51527ac46151c56c766b52527ac46c766b00c310526567496442795075626c69634b6579876c766b53527ac46c766b53c3646d00616c766b51c3c0529c009c6c766b54527ac46c766b54c3641000616c766b52c30052c461623700616c766b51c300c36c766b55527ac46c766b51c351c36c766b56527ac46c766b52c3006c766b55c36c766b56c3617c65740ac4616c766b52c36c766b57527ac462d0066c766b00c3055265674964876c766b58527ac46c766b58c3648000616c766b51c3c0539c009c6c766b59527ac46c766b59c3641b00616c766b52c30052c46c766b52c36c766b57527ac4628306616c766b51c300c36c766b5a527ac46c766b51c351c36c766b5b527ac46c766b52c3006c766b5ac36c766b5bc36c766b51c352c361527265630bc4616c766b52c36c766b57527ac46238066c766b00c3064164644b6579876c766b5c527ac46c766b5cc3648100616c766b51c3c0539c009c6c766b5d527ac46c766b5dc3641000616c766b52c30052c461624b00616c766b51c300c36c766b5e527ac46c766b51c351c36c766b5f527ac46c766b51c352c36c766b60527ac46c766b52c3006c766b5ec36c766b5fc36c766b60c3615272651611c4616c766b52c36c766b57527ac4629e056c766b00c30952656d6f76654b6579876c766b0111527ac46c766b0111c3649400616c766b51c3c0539c009c6c766b0112527ac46c766b0112c3641b00616c766b52c30052c46c766b52c36c766b57527ac4624905616c766b51c300c36c766b0113527ac46c766b51c351c36c766b0114527ac46c766b51c352c36c766b0115527ac46c766b52c3006c766b0113c36c766b0114c36c766b0115c361527265f511c4616c766b52c36c766b57527ac462ec046c766b00c30c416464417474726962757465876c766b0116527ac46c766b0116c364d000616c766b51c3c0559c009c6c766b0117527ac46c766b0117c3641b00616c766b52c30052c46c766b52c36c766b57527ac4629404616c766b51c300c36c766b0118527ac46c766b51c351c36c766b0119527ac46c766b51c352c36c766b011a527ac46c766b51c353c36c766b011b527ac46c766b51c354c36c766b011c527ac46c766b52c3006c766b0118c36c766b0119c36c766b011ac36c766b011bc36c766b011cc3615479517956727551727553795279557275527275650d15c4616c766b52c36c766b57527ac462fb036c766b00c30b4164645265636f76657279876c766b011d527ac46c766b011dc3648900616c766b51c3c0539c009c6c766b011e527ac46c766b011ec3641000616c766b52c30052c461625100616c766b51c300c36c766b011f527ac46c766b51c351c36c766b0120527ac46c766b51c352c36c766b0121527ac46c766b52c3006c766b011fc36c766b0120c36c766b0121c3615272652712c4616c766b52c36c766b57527ac46252036c766b00c30e4368616e67655265636f76657279876c766b0122527ac46c766b0122c3648900616c766b51c3c0539c009c6c766b0123527ac46c766b0123c3641000616c766b52c30052c461625100616c766b51c300c36c766b0124527ac46c766b51c351c36c766b0125527ac46c766b51c352c36c766b0126527ac46c766b52c3006c766b0124c36c766b0125c36c766b0126c361527265c312c4616c766b52c36c766b57527ac462a6026c766b00c3114164644174747269627574654172726179876c766b0127527ac46c766b0127c3641b00616c766b52c30000c46c766b52c36c766b57527ac46265026c766b00c30f52656d6f7665417474726962757465876c766b0128527ac46c766b0128c3648100616c766b51c3c0529c009c6c766b0129527ac46c766b0129c36408006161625100616c766b51c300c36c766b012a527ac46c766b51c351c36c766b012b527ac46c766b51c352c36c766b012c527ac46c766b52c3006c766b012ac36c766b012bc36c766b012cc361527265e214c4616c766b52c36c766b57527ac462c0016c766b00c30d4765745075626c69634b657973876c766b012d527ac46c766b012dc3645d00616c766b51c3c0519c009c6c766b012e527ac46c766b012ec3641000616c766b52c30052c461622500616c766b51c300c36c766b012f527ac46c766b52c3006c766b012fc361652521c4616c766b52c36c766b57527ac46241016c766b00c30647657444444f876c766b0130527ac46c766b0130c3648c00616c766b51c3c0519c009c6c766b0131527ac46c766b0131c3641000616c766b52c30052c461625400616c766b51c300c36c766b0132527ac46c766b0132c36165b3206c766b0133527ac46c766b0132c36165ff216c766b0134527ac46c766b52c3006c766b0133c3616532246c766b0134c3616528247ec4616c766b52c36c766b57527ac4629a006c766b00c30d47657441747472696275746573876c766b0135527ac46c766b0135c3645d00616c766b51c3c0519c009c6c766b0136527ac46c766b0136c3641000616c766b52c30052c461622500616c766b51c300c36c766b0137527ac46c766b52c3006c766b0137c361655d21c4616c766b52c36c766b57527ac4621b00616c766b52c30000c46c766b52c36c766b57527ac46203006c766b57c3616c756654c56b6c766b00527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c3617c680f4e656f2e53746f726167652e4765746c766b51527ac46c766b51c3c0641b006c766b51c300517f5100517f907c907c9e63070051620400006c766b52527ac46c766b52c36c766b53527ac46203006c766b53c3616c756651c56b6c766b00527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c351615272680f4e656f2e53746f726167652e50757461616c756653c56b6c766b00527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c3547e617c680f4e656f2e53746f726167652e4765746c766b51527ac46c766b51c36c766b52527ac46203006c766b52c3616c756652c56b6c766b00527ac46c766b51527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c3547e6c766b51c3615272680f4e656f2e53746f726167652e50757461616c756653c56b6c766b00527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c3517e617c680f4e656f2e53746f726167652e4765746c766b51527ac46c766b51c36c766b52527ac46203006c766b52c3616c756652c56b6c766b00527ac46c766b51527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c3517e6c766b51c3615272680f4e656f2e53746f726167652e50757461616c756652c56b6c766b00527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c3577e617c680f4e656f2e53746f726167652e4765746c766b51527ac46203006c766b51c3616c756652c56b6c766b00527ac46c766b51527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c3577e6c766b51c3615272680f4e656f2e53746f726167652e50757461616c756653c56b6c766b00527ac4616c766b00c3c06410006c766b00c3c002ff00a0620400516c766b51527ac46c766b51c3640e00006c766b52527ac4621d006c766b00c3c06165911b6c766b00c37e6c766b52527ac46203006c766b52c3616c756653c56b6c766b00527ac4616c766b00c3c0519f6319006c766b00c3c06c766b00c300517f51939c009c620400516c766b51527ac46c766b51c3640e00006c766b52527ac4621c006c766b00c3516c766b00c300517f7f6c766b52527ac46203006c766b52c3616c756659c56b6c766b00527ac46c766b51527ac4616100c56c766b52527ac46c766b00c3616516ff6c766b53527ac46c766b53c3c0519f6c766b54527ac46c766b54c3640f0061526c766b55527ac46232016c766b53c3616521fc6c766b56527ac46c766b56c3641b0061566c766b52527ac46c766b52c36c766b55527ac46202016c766b51c36168184e656f2e52756e74696d652e436865636b5769746e657373009c6c766b57527ac46c766b57c3643d00610e696e76616c69642063616c6c657261680f4e656f2e52756e74696d652e4c6f6761536c766b52527ac46c766b52c36c766b55527ac46297006c766b53c36c766b51c3617c655e106c766b58527ac46c766b58c3646200616c766b53c36165ebfb616c766b53c300617c657ffc616c766b53c351617c6520fd610872656769737465726c766b00c3617c08526567697374657253c168124e656f2e52756e74696d652e4e6f7469667961516c766b52527ac461620b00546c766b52527ac46c766b52c36c766b55527ac46203006c766b55c3616c75660121c56b6c766b00527ac46c766b51527ac46c766b52527ac4616100c56c766b53527ac46c766b00c3616587fd6c766b54527ac46c766b54c3c0519f6c766b5c527ac46c766b5cc3640e00526c766b5d527ac462f1056c766b54c3616593fa6c766b5e527ac46c766b5ec3641b0061566c766b53527ac46c766b53c36c766b5d527ac462c1056c766b51c36168184e656f2e52756e74696d652e436865636b5769746e657373009c6c766b5f527ac46c766b5fc3641b0061536c766b53527ac46c766b53c36c766b5d527ac4627805006c766b55527ac4516c766b56527ac4006c766b57527ac46c766b52c3640e006c766b52c300517f620400006c766b58527ac46c766b52c3c051946c766b59527ac4020001c56c766b5a527ac4006c766b55527ac4516c766b56527ac4627304616c766b59c3529f6c766b0118527ac46c766b0118c3644f0061526c766b53527ac412496e76616c6964204174747269627574657361680f4e656f2e52756e74696d652e4c6f6761516c766b57c3966c766b57527ac46c766b53c36c766b5d527ac462b3046c766b52c36c766b56c3517f020001956c766b52c36c766b56c35193517f936c766b60527ac46c766b59c36c766b60c352939f6c766b0119527ac46c766b0119c3644f0061526c766b53527ac412496e76616c6964204174747269627574657361680f4e656f2e52756e74696d652e4c6f6761516c766b57c3966c766b57527ac46c766b53c36c766b5d527ac46223046c766b52c36c766b56c352936c766b60c37f6c766b0111527ac46c766b0111c36411006c766b0111c3c051a0009c620400516c766b011a527ac46c766b011ac3644f0061526c766b53527ac412496e76616c6964204174747269627574657361680f4e656f2e52756e74696d652e4c6f6761516c766b57c3966c766b57527ac46c766b53c36c766b5d527ac46294036c766b0111c300517f6c766b0112527ac46c766b0111c3c0516c766b0112c393a0009c6c766b011b527ac46c766b011bc3644f0061526c766b53527ac412496e76616c6964204174747269627574657361680f4e656f2e52756e74696d652e4c6f6761516c766b57c3966c766b57527ac46c766b53c36c766b5d527ac46214036c766b0111c3516c766b0112c393517f6c766b0113527ac46c766b0111c3c0546c766b0112c3936c766b0113c3939f6c766b011c527ac46c766b011cc3644f0061526c766b53527ac412496e76616c6964204174747269627574657361680f4e656f2e52756e74696d652e4c6f6761516c766b57c3966c766b57527ac46c766b53c36c766b5d527ac46288026c766b0111c3526c766b0112c3936c766b0113c393517f020001956c766b0111c3536c766b0112c3936c766b0113c393517f936c766b0114527ac46c766b0112c36c766b0113c3936c766b0114c39354936c766b60c39c009c6c766b011d527ac46c766b011dc3644f0061526c766b53527ac412496e76616c6964204174747269627574657361680f4e656f2e52756e74696d652e4c6f6761516c766b57c3966c766b57527ac46c766b53c36c766b5d527ac462d2016c766b0111c3516c766b0112c37f6c766b0115527ac46c766b0111c3526c766b0112c3936c766b0113c37f6c766b0116527ac46c766b0111c3546c766b0112c3936c766b0113c3936c766b0114c37f6c766b0117527ac46c766b0113c302ff00a06c766b011e527ac46c766b011ec3641b0061526c766b53527ac46c766b53c36c766b5d527ac46248016c766b54c36c766b0115c3617c652d0b6c766b011f527ac46c766b011fc3640800616162050061616c766b54c36c766b0115c36c766b0116c36c766b0117c3615379517955727551727552795279547275527275652213616c766b5ac36c766b55c36c766b0115c3c46c766b56c352936c766b60c3936c766b56527ac46c766b59c3526c766b60c393946c766b59527ac4616c766b55c351936c766b55527ac46c766b55c36c766b58c39f6c766b0120527ac46c766b0120c36377fb6c766b55c36c766b5b527ac46c766b54c36165a6f5616c766b54c36c766b51c3617c65f009756c766b54c351617c65d7f6616c766b54c36c766b5bc3617c651af6610872656769737465726c766b00c3617c08526567697374657253c168124e656f2e52756e74696d652e4e6f7469667961516c766b53527ac46c766b53c36c766b5d527ac46203006c766b5dc3616c75665bc56b6c766b00527ac46c766b51527ac46c766b52527ac4616100c56c766b53527ac46c766b00c361653bf76c766b54527ac46c766b54c3c0519f6c766b56527ac46c766b56c3640e00526c766b57527ac46236016c766b54c3616547f4009c6c766b58527ac46c766b58c3641b0061556c766b53527ac46c766b53c36c766b57527ac46204016c766b54c36c766b52c3617c6598086428006c766b52c36168184e656f2e52756e74696d652e436865636b5769746e657373009c620400516c766b59527ac46c766b59c3641b0061536c766b53527ac46c766b53c36c766b57527ac462a5006c766b54c361652af56c766b55527ac46c766b54c36c766b51c3617c657e086c766b5a527ac46c766b5ac3645f00616c766b54c36c766b55c35193617c6550f561036164646c766b00c36c766b51c3615272095075626c69634b657954c168124e656f2e52756e74696d652e4e6f7469667961516c766b53527ac46c766b53c36c766b57527ac4621b0061546c766b53527ac46c766b53c36c766b57527ac46203006c766b57c3616c75665bc56b6c766b00527ac46c766b51527ac46c766b52527ac4616100c56c766b53527ac46c766b00c36165aaf56c766b54527ac46c766b54c3c0519f6c766b56527ac46c766b56c3640e00526c766b57527ac46239016c766b54c36165b6f2009c6c766b58527ac46c766b58c3641b0061556c766b53527ac46c766b53c36c766b57527ac46207016c766b54c36c766b52c3617c6507076428006c766b52c36168184e656f2e52756e74696d652e436865636b5769746e657373009c620400516c766b59527ac46c766b59c3641b0061536c766b53527ac46c766b53c36c766b57527ac462a8006c766b54c3616599f36c766b55527ac46c766b54c36c766b51c3617c6523076c766b5a527ac46c766b5ac3646200616c766b54c36c766b55c35194617c65bff3610672656d6f76656c766b00c36c766b51c3615272095075626c69634b657954c168124e656f2e52756e74696d652e4e6f7469667961516c766b53527ac46c766b53c36c766b57527ac4621b0061586c766b53527ac46c766b53c36c766b57527ac46203006c766b57c3616c756653c56b6c766b00527ac46c766b51527ac4616c766b51c36c766b00c3616581f3617c65a90f6c766b52527ac46203006c766b52c3616c75665bc56b6c766b00527ac46c766b51527ac46c766b52527ac4616100c56c766b53527ac46c766b00c36165def36c766b54527ac46c766b54c3c0519f6c766b56527ac46c766b56c3640e00526c766b57527ac462ed006c766b54c36165eaf0009c6c766b58527ac46c766b58c3641b0061556c766b53527ac46c766b53c36c766b57527ac462bb006c766b54c36c766b52c3617c653b056428006c766b52c36168184e656f2e52756e74696d652e436865636b5769746e657373009c620400516c766b59527ac46c766b59c3641b0061536c766b53527ac46c766b53c36c766b57527ac4625c006c766b54c361657af26c766b55527ac46c766b55c300a06c766b5a527ac46c766b5ac3640f0061006c766b57527ac4622a006c766b54c36c766b51c3617c6593f261516c766b53527ac46c766b53c36c766b57527ac46203006c766b57c3616c756659c56b6c766b00527ac46c766b51527ac46c766b52527ac4616100c56c766b53527ac46c766b00c3616596f26c766b54527ac46c766b54c3c0519f6c766b55527ac46c766b55c3640e00526c766b56527ac4629a006c766b54c36165a2ef009c6c766b57527ac46c766b57c3640e00006c766b56527ac46275006c766b52c36c766b54c3616599f1617c65c10d6428006c766b52c36168184e656f2e52756e74696d652e436865636b5769746e657373009c620400516c766b58527ac46c766b58c3640f0061006c766b56527ac4621e006c766b54c36c766b51c3617c6592f161516c766b56527ac46203006c766b56c3616c75665dc56b6c766b00527ac46c766b51527ac46c766b52527ac46c766b53527ac46c766b54527ac4616100c56c766b55527ac46c766b00c3616593f16c766b56527ac46c766b56c3c0519f6c766b57527ac46c766b57c3640e00526c766b58527ac462a7016c766b56c361659fee009c6c766b59527ac46c766b59c3641b0061556c766b55527ac46c766b55c36c766b58527ac46275016c766b56c36c766b54c3617c65f0026428006c766b54c36168184e656f2e52756e74696d652e436865636b5769746e657373009c620400516c766b5a527ac46c766b5ac3641b0061536c766b55527ac46c766b55c36c766b58527ac46216016c766b56c36c766b51c3617c6552036c766b5b527ac46c766b5bc3648900616c766b56c36165b6ee6c766b5c527ac46c766b56c36c766b5cc35193617c65fbee616c766b56c36c766b51c36c766b52c36c766b53c361537951795572755172755279527954727552727565300b61036164646c766b00c36c766b51c36152720941747472696275746554c168124e656f2e52756e74696d652e4e6f746966796161626700616c766b56c36c766b51c36c766b52c36c766b53c361537951795572755172755279527954727552727565cc0a61067570646174656c766b00c36c766b51c36152720941747472696275746554c168124e656f2e52756e74696d652e4e6f746966796161516c766b58527ac46203006c766b58c3616c75665bc56b6c766b00527ac46c766b51527ac46c766b52527ac4616100c56c766b53527ac46c766b00c3616591ef6c766b54527ac46c766b54c3c0519f6c766b55527ac46c766b55c3640e00526c766b56527ac46226016c766b54c361659dec009c6c766b57527ac46c766b57c3641b0061556c766b53527ac46c766b53c36c766b56527ac462f4006c766b54c36c766b52c3617c65ee006428006c766b52c36168184e656f2e52756e74696d652e436865636b5769746e657373009c620400516c766b58527ac46c766b58c3641b0061536c766b53527ac46c766b53c36c766b56527ac46295006c766b54c36c766b51c3617c6586016c766b59527ac46c766b59c3646c00616c766b54c36165b4ec6c766b5a527ac46c766b54c36c766b5ac35194617c65f9ec616c766b54c36c766b51c3617c65c009610672656d6f76656c766b00c36c766b51c36152720941747472696275746554c168124e656f2e52756e74696d652e4e6f746966796161516c766b56527ac46203006c766b56c3616c756654c56b6c766b00527ac46c766b51527ac4616c766b00c3526c766b51c3615272655f07009c6c766b52527ac46c766b52c3640f0061006c766b53527ac4620f0061516c766b53527ac46203006c766b53c3616c756653c56b6c766b00527ac46c766b51527ac4616c766b00c3526c766b51c361527265b8006c766b52527ac46203006c766b52c3616c756653c56b6c766b00527ac46c766b51527ac4616c766b00c3526c766b51c3615272654a026c766b52527ac46203006c766b52c3616c756653c56b6c766b00527ac46c766b51527ac4616c766b00c3556c766b51c3615272654c006c766b52527ac46203006c766b52c3616c756653c56b6c766b00527ac46c766b51527ac4616c766b00c3556c766b51c361527265de016c766b52527ac46203006c766b52c3616c756659c56b6c766b00527ac46c766b51527ac46c766b52527ac4616c766b00c36c766b51c3617c6596066c766b53527ac46c766b00c36c766b51c36c766b52c36152726511066c766b54527ac46c766b00c36c766b51c36c766b53c361527265f5056c766b55527ac46c766b54c300a06c766b56527ac46c766b56c3640f0061006c766b57527ac46239016c766b53c3009c6c766b58527ac46c766b58c3646400616168164e656f2e53746f726167652e476574436f6e746578746c766b00c36c766b51c36c766b52c37e7e0000617c65fd04615272680f4e656f2e53746f726167652e507574616c766b00c36c766b51c36c766b52c3615272651f06616162b700616168164e656f2e53746f726167652e476574436f6e746578746c766b00c36c766b51c36c766b52c37e7e6c766b53c300617c659804615272680f4e656f2e53746f726167652e507574616168164e656f2e53746f726167652e476574436f6e746578746c766b00c36c766b51c36c766b53c37e7e6c766b55c3616557036c766b52c3617c654604615272680f4e656f2e53746f726167652e507574616c766b00c36c766b51c36c766b52c36152726568056161516c766b57527ac46203006c766b57c3616c756660c56b6c766b00527ac46c766b51527ac46c766b52527ac4616c766b00c36c766b51c36c766b52c3615272655f046c766b53527ac46c766b53c3009c6c766b57527ac46c766b57c3640f0061006c766b58527ac462ad02006c766b54527ac46c766b53c36165f7026c766b55527ac46c766b53c3616595026c766b56527ac46c766b55c3009c6c766b59527ac46c766b59c364bb00616c766b56c3009c6c766b5a527ac46c766b5ac3641e00616c766b00c36c766b51c36c766b54c36152726590046161628500616c766b00c36c766b51c36c766b56c361527265b0036c766b5b527ac46168164e656f2e53746f726167652e476574436f6e746578746c766b00c36c766b51c36c766b56c37e7e6c766b5bc36165f60100617c65e902615272680f4e656f2e53746f726167652e507574616c766b00c36c766b51c36c766b56c3615272650b04616161626801616c766b56c3009c6c766b5c527ac46c766b5cc3647200616c766b00c36c766b51c36c766b55c36152726513036c766b5d527ac46168164e656f2e53746f726167652e476574436f6e746578746c766b00c36c766b51c36c766b55c37e7e006c766b5dc36165aa01617c654c02615272680f4e656f2e53746f726167652e507574616162e100616c766b00c36c766b51c36c766b56c361527265a4026c766b5e527ac46c766b00c36c766b51c36c766b55c36152726588026c766b5f527ac46168164e656f2e53746f726167652e476574436f6e746578746c766b00c36c766b51c36c766b55c37e7e6c766b56c36c766b5fc361651b01617c65bd01615272680f4e656f2e53746f726167652e507574616168164e656f2e53746f726167652e476574436f6e746578746c766b00c36c766b51c36c766b56c37e7e6c766b5ec361657c006c766b55c3617c656b01615272680f4e656f2e53746f726167652e5075746161616168164e656f2e53746f726167652e476574436f6e746578746c766b00c36c766b51c36c766b52c37e7e6c766b54c3615272680f4e656f2e53746f726167652e50757461516c766b58527ac46203006c766b58c3616c756654c56b6c766b00527ac4616c766b00c300517f6c766b51527ac46c766b00c3516c766b51c37f6c766b52527ac46c766b00c3516c766b51c3936c766b52c37f6c766b53527ac46203006c766b53c3616c756657c56b6c766b00527ac4616c766b00c300517f6c766b51527ac46c766b00c3516c766b51c37f6c766b52527ac46c766b00c3526c766b51c3936c766b52c393517f6c766b53527ac46c766b53c300517f6c766b54527ac46c766b00c3536c766b51c3936c766b52c3936c766b54c37f6c766b55527ac46c766b00c3536c766b51c3936c766b52c3936c766b54c3936c766b55c37f6c766b56527ac46203006c766b56c3616c756657c56b6c766b00527ac46c766b51527ac4616c766b00c3c06c766b54527ac46c766b51c3c06c766b55527ac46c766b54c3c06c766b52527ac46c766b55c3c06c766b53527ac46c766b52c36165c9026c766b54c36c766b00c301016c766b53c36165b4026c766b55c36c766b51c37e7e7e7e7e7e6c766b56527ac46203006c766b56c3616c756654c56b6c766b00527ac46c766b51527ac46c766b52527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c36c766b51c36c766b52c37e7e617c680f4e656f2e53746f726167652e4765746c766b53527ac46203006c766b53c3616c756653c56b6c766b00527ac46c766b51527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c36c766b51c37e617c680f4e656f2e53746f726167652e4765746c766b52527ac46203006c766b52c3616c756653c56b6c766b00527ac46c766b51527ac46c766b52527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c36c766b51c37e6c766b52c3615272680f4e656f2e53746f726167652e50757461616c756653c56b6c766b00527ac46c766b51527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c3566c766b51c37e7e617c680f4e656f2e53746f726167652e4765746c766b52527ac46203006c766b52c3616c756654c56b6c766b00527ac46c766b51527ac46c766b52527ac46c766b53527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c3566c766b51c37e7e6c766b52c3c06165c3006c766b52c36c766b53c37e7e615272680f4e656f2e53746f726167652e50757461616c756652c56b6c766b00527ac46c766b51527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c3566c766b51c37e7e617c68124e656f2e53746f726167652e44656c65746561616c756653c56b6c766b00527ac46c766b51527ac4616c766b00c3c06c766b51c3c0907c907c9e6311006c766b00c36c766b51c39c620400006c766b52527ac46203006c766b52c3616c756653c56b6c766b00527ac4614d0001000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c9cacbcccdcecfd0d1d2d3d4d5d6d7d8d9dadbdcdddedfe0e1e2e3e4e5e6e7e8e9eaebecedeeeff0f1f2f3f4f5f6f7f8f9fafbfcfdfeff6c766b51527ac46c766b51c36c766b00c3517f6c766b52527ac46203006c766b52c3616c75665bc56b6c766b00527ac4616c766b00c36165e7e26c766b51527ac400006c766b53527ac46168164e656f2e53746f726167652e476574436f6e746578746c766b51c3527e617c680f4e656f2e53746f726167652e4765746c766b54527ac46c766b54c3c0009c6c766b55527ac46c766b55c3640f0061006c766b56527ac462d700616c766b54c36c766b57527ac46c766b57c3616515036c766b52527ac4629300616168164e656f2e53746f726167652e476574436f6e746578746c766b51c3527e6c766b57c37e617c680f4e656f2e53746f726167652e4765746c766b58527ac46c766b58c36165d6f96c766b57527ac46c766b53c351936c766b53527ac46c766b57c3c0009c6c766b59527ac46c766b59c36406006225006c766b52c36c766b57c3616584027e6c766b52527ac461516c766b5a527ac46268ff6c766b53c361658afd6c766b52c37e6c766b56527ac46203006c766b56c3616c75665cc56b6c766b00527ac4616c766b00c3616589e16c766b51527ac400006c766b53527ac46168164e656f2e53746f726167652e476574436f6e746578746c766b51c3557e617c680f4e656f2e53746f726167652e4765746c766b54527ac46c766b54c3c0009c6c766b55527ac46c766b55c3640f0061006c766b56527ac4627301616c766b54c36c766b57527ac46168164e656f2e53746f726167652e476574436f6e746578746c766b51c3567e6c766b57c37e617c680f4e656f2e53746f726167652e4765746c766b58527ac46c766b57c3616577016c766b58c361656e017e616569016c766b52527ac462e100616168164e656f2e53746f726167652e476574436f6e746578746c766b51c3557e6c766b57c37e617c680f4e656f2e53746f726167652e4765746c766b59527ac46c766b59c361652af86c766b57527ac46c766b53c351936c766b53527ac46c766b57c3c0009c6c766b5a527ac46c766b5ac36406006273006168164e656f2e53746f726167652e476574436f6e746578746c766b51c3567e6c766b57c37e617c680f4e656f2e53746f726167652e4765746c766b58527ac46c766b52c36c766b57c3616598006c766b58c361658f007e61658a007e6c766b52527ac461516c766b5b527ac4621aff6c766b53c3616590fb6c766b52c37e6c766b56527ac46203006c766b56c3616c756654c56b6c766b00527ac4616c766b00c3616597fc6c766b51527ac46c766b00c36165e5fd6c766b52527ac46c766b51c3616520006c766b52c3616517007e6c766b53527ac46203006c766b53c3616c756657c56b6c766b00527ac4616c766b00c3c06c766b51527ac46c766b51c3020001976c766b52527ac46c766b51c36c766b52c394020001966c766b51527ac46c766b51c3020001976c766b53527ac46c766b51c36c766b53c394020001966c766b51527ac46c766b51c3020001976c766b54527ac46c766b51c36c766b54c394020001966c766b51527ac46c766b51c3020001976c766b55527ac46c766b55c3616583fa6c766b54c361657afa7e6c766b53c3616570fa7e6c766b52c3616566fa7e6c766b00c37e6c766b56527ac46203006c766b56c3616c756601046e616d6503312e3001310131013101002466386537653361662d323463322d343135352d626361382d66313866636463613834653400000000000000000000")

	var txn types.Transaction
	if err := txn.Deserialize(bytes.NewReader(bys)); err != nil {
		fmt.Print("Deserialize Err:",err)
		os.Exit(0)
	}
	fmt.Printf("TxType:%x\n", txn.TxType)
	os.Exit(0)
}