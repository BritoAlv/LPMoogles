#include <bits/stdc++.h>
#include <httpserver.hpp>
#include <filesystem>
#include <fstream>
using namespace httpserver;
using namespace std;

vector<string> splitInWords(string &text)
{
    // provide better implementation.
    auto isletter = [](char c) { return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z'); };
    vector<string> words = vector<string>();
    for (int i = 0; i < text.size(); i++)
    {
        if(!isletter(text[i])){
            continue;
        }
        int start = i;
        int end = i;
        if (isletter(text[i]))
        {
            while (end + 1 < text.size() && isletter(text[end + 1]))
            {
                end++;
            }
        }
        string w = "";
        for (int j = start; j <= end; j++)
        {
            w += text[j];
        }
        words.push_back(w);
        i = end;
    }
    return words;
}

struct ResultToWebDto
{
    string Name;
    string Snippet;
};

struct ResultFromDto
{
    string Name;
    string Text;
};

ResultToWebDto Wrapper(ResultFromDto &o)
{
    ResultToWebDto A = *(new ResultToWebDto);
    A.Name = o.Name;
    A.Snippet = "";
    for (int i = 0; i < min(100, (int)o.Text.size()); i++)
    {
        A.Snippet += o.Text[i];
    }
    return A;
}

class ModelTfIdf
{
    public:
        map<string, int> Idf;
        vector<ResultFromDto> Items;
        map<string, map<string, int>> Tf;
        unsigned long long TotalDocuments;
        ModelTfIdf(vector<ResultFromDto> &arr)
        {
            map<string, int> Idf = *(new map<string, int>());
            map<string, map<string, int>> Tf = *(new map<string, map<string, int>>());
            this->TotalDocuments = arr.size();
            for (ResultFromDto &doc : arr)
            {
                map<string, int> tfmap = *(new map<string, int>());
                auto wordsInDoc = splitInWords((doc.Text));
                for (string &word : wordsInDoc)
                {
                    if (Idf.find(word) == Idf.end())
                    {
                        Idf[word] = 1;
                    }
                    if (tfmap.find(word) == tfmap.end())
                    {
                        tfmap[word] = 0;
                    }
                    tfmap[word]++;
                }
                Tf[doc.Name] = tfmap;
            }
            this->Idf = Idf;
            this->Tf = Tf;
            this->Items = arr;    
        }
};


vector<float> tdIdfCalculator(vector<string> &words, map<string, int> &source, ModelTfIdf *model)
{
    auto result = vector<float>(words.size(), 0);
    for (int index = 0; index < words.size(); index++)
    {
        auto word = words[index];
        auto v = source.find(word);
        if (v != source.end())
        {
            float v1 = (*v).second;
            v1 *= log((float)model->TotalDocuments / (float)model->Idf[word] + 1);
            result[index] = v1;
        }
    }
    return result;
}

map<string, int> QueryWords(vector<string> &wordsInDoc, ModelTfIdf &model)
{
    map<string, int> tfmap = *(new map<string, int>());
    for (auto &word : wordsInDoc)
    {
        if (tfmap.find(word) == tfmap.end())
        {
            tfmap[word] = 0;
        }
        tfmap[word]++;
    }
    return tfmap;
}

float cos_sim(vector<float> A, vector<float> B)
{
    float result = 0;
    float mag_A = 0;
    float mag_B = 0;
    for (int i = 0; i < A.size(); i++)
    {
        result += A[i] * B[i];
        mag_A += A[i] * A[i];
        mag_B += B[i] * B[i];
    }
    if (mag_A == 0 || mag_B == 0)
    {
        return 0;
    }
    return result / (sqrt(mag_A) * sqrt(mag_B));
}

vector<ResultFromDto> read_txt_files_local(const std::string& folderPath)
{
    vector<ResultFromDto> ans = vector<ResultFromDto>();
    for (const auto& entry : std::filesystem::directory_iterator(folderPath)) {
        if (entry.path().extension() == ".txt") {
            
            std::ifstream inFile(entry.path());
            std::string line;
            ResultFromDto current = *(new ResultFromDto());
            current.Name = entry.path().filename();
            while (std::getline(inFile, line)) {
                current.Text += line;
            }
            ans.push_back(current);
        }
    }
    return ans;
}

ModelTfIdf model;

void setup_model()
{
    vector<ResultFromDto> items = read_txt_files_local("./database");
    model = *(new ModelTfIdf(items));
    return;
}

vector<ResultToWebDto> startSearchFromQuery(string inputValue)
{
    auto querywords = splitInWords(inputValue);
    auto querySource = QueryWords(querywords, model);
    auto queryTfIdf = tdIdfCalculator(querywords, querySource, &model);

    vector<pair<float, int>> docs_values;
    for (int i = 0; i < model.TotalDocuments; i++)
    {
        docs_values.push_back(
            {cos_sim(tdIdfCalculator(querywords, model.Tf[model.Items[i].Name], &model), queryTfIdf), i});
    }
    sort(docs_values.begin(), docs_values.end());
    vector<ResultToWebDto> ans;
    for (int i = model.TotalDocuments - 1; i >= 0; i--)
    {
        ans.push_back(Wrapper(model.Items[docs_values[i].second]));
    }
    return ans;
}

class hello_world_resource : public http_resource
{
  public:
    shared_ptr<http_response> render_GET(const http_request &req);
    void set_some_data(const string &s)
    {
        data = s;
    }
    string data;
};

shared_ptr<http_response> hello_world_resource::render_GET(const http_request &req)
{
    string_view datapar = req.get_arg("query");
    set_some_data(datapar == "" ? "Waiting For A Query" : string(datapar));
    vector<ResultToWebDto> to_show = startSearchFromQuery(this->data);
    string result = "";
    for (auto &x : to_show)
    {
        result += x.Name;
        result += '\n';
        result += x.Snippet;
        result += '\n';
    }
    return shared_ptr<http_response>(new string_response(result, 200));
}

int main()
{
    setup_model();
    webserver ws = create_webserver(8080);
    hello_world_resource hwr;
    ws.register_resource("/", &hwr, true);
    ws.start(true);
    return 0;
}