import joblib
import json
import pandas as pd
import sys

required_features = ['umur', 'monthly_income', 'payroll', 'gender_MALE', 'marital_status_Single', 'transaction_activity_Inactive', 'categorysegmen_BUMN', 'categorysegmen_Lembaga Negara', 'categorysegmen_Non Target Market', 'categorysegmen_Pendidikan', 'categorysegmen_Pensiun', 'categorysegmen_RS', 'categorysegmen_Swasta']
produk_list=['mitraguna', 'hasanahcard', 'griya', 'oto', 'pensiun', 'prapensiun']
scaler = joblib.load('scripts/scaler.pkl')

def predict_final_deploy(data_user):
    input_dict = {
        'umur': data_user['umur'],
        'monthly_income': data_user['income'],
        'payroll': data_user['payroll'],
        'marital_status_Single': 1 if data_user['marital_status'] == 0 else 0,
        'transaction_activity_Inactive': 1 if data_user['transaction_activity'].lower() == 'inactive' else 0,
        'gender_MALE': 1 if data_user['gender'].lower() == 'male' else 0
    }

    category_segmen_values = ['BO2', 'Swasta', 'Pendidikan', 'BUMN', 'Non Target Market', 'Lembaga Negara', 'Pensiun', 'RS']
    for segmen in category_segmen_values:
        colname = f'categorysegmen_{segmen.lower()}'
        input_dict[colname] = 1 if data_user['category_segmen'] == segmen else 0

    existing_products = data_user['existing_product']
    if not isinstance(existing_products, list):
        existing_products = [existing_products]

    input_data = pd.DataFrame([input_dict])
    for col in required_features:
        if col not in input_data.columns:
            input_data[col] = 0

    input_data = input_data[required_features]
    numerical_cols = ['umur', 'monthly_income']
    input_data[numerical_cols] = scaler.transform(input_data[numerical_cols])

    pred_proba = {}
    for produk in produk_list:
        model_filename = f'scripts/model_{produk}.pkl'
        model = joblib.load(model_filename)
        pred_prob = model.predict_proba(input_data)[:, 1][0]
        pred_proba[produk] = pred_prob

    pred_df = pd.DataFrame([pred_proba])

    for produk in existing_products:
        if produk in pred_df.columns:
            pred_df[produk] = 0

    if data_user['payroll'] == 0:
        if 'mitraguna' in pred_df.columns:
            pred_df['mitraguna'] = 0

    return pred_df

if __name__ == "__main__":
    input_json = sys.stdin.read()
    dummy_data = json.loads(input_json)

    dummy_pred_df = predict_final_deploy(dummy_data)

    top_preds = dummy_pred_df.T.sort_values(by=0, ascending=False).head(3)[0].to_dict()

    print(json.dumps(top_preds))